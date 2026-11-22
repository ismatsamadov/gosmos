package main

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"sync"
	"time"
)

// keyID names one customer's idempotency key: Brinewatch's k1 is not Cirrine's k1.
type keyID struct {
	Customer string
	Key      string
}

// kept is a request that carried an idempotency key: its body's fingerprint and, once it is done, its answer.
type kept struct {
	Sum    [sha256.Size]byte // SHA-256 of the request's body
	At     time.Time         // when the key was claimed
	Done   bool              // false while the first request is still running
	Status int
	Header http.Header // the answer's Content-Type and Location
	Body   []byte
}

// keyStore keeps idempotency keys and their answers for 24 hours. It is safe for many requests at once.
type keyStore struct {
	now  func() time.Time
	mu   sync.Mutex
	keys map[keyID]*kept
}

// newKeyStore returns an empty keyStore on the clock now.
func newKeyStore(now func() time.Time) *keyStore {
	return &keyStore{now: now, keys: map[keyID]*kept{}}
}

// claim decides what a request with key id and this body gets: a refusal, a kept answer to replay, or nil, nil to run it.
func (s *keyStore) claim(id keyID, body []byte) (*kept, error) {
	sum := sha256.Sum256(body)
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	if k, ok := s.keys[id]; ok && now.Sub(k.At) < 24*time.Hour {
		switch {
		case k.Sum != sum:
			return nil, differentRequest(id.Key)
		case !k.Done:
			return nil, inProgress(id.Key)
		}
		return k, nil
	}
	s.keys[id] = &kept{Sum: sum, At: now}
	return nil, nil
}

// differentRequest is the 409 for a key that comes back with another body.
func differentRequest(key string) error {
	return &APIError{Status: http.StatusConflict, Code: "conflict", Message: fmt.Sprintf("idempotency key %q was used for a different request", key)}
}

// inProgress is the 409 for a key whose first request is still running.
func inProgress(key string) error {
	return &APIError{Status: http.StatusConflict, Code: "conflict", Message: fmt.Sprintf("a request with idempotency key %q is in progress", key)}
}

// finish keeps the answer to a claimed request, or forgets the key when the answer is 500 or above.
func (s *keyStore) finish(id keyID, status int, header http.Header, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k, ok := s.keys[id]
	if !ok || k.Done {
		return
	}
	if status >= 500 {
		delete(s.keys, id)
		return
	}
	k.Status, k.Header, k.Body, k.Done = status, header, body, true
}

// answerRecorder passes an answer on to the caller and keeps a copy of its status and body.
type answerRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

// WriteHeader records the first status and passes it on.
func (a *answerRecorder) WriteHeader(code int) {
	if a.status == 0 {
		a.status = code
	}
	a.ResponseWriter.WriteHeader(code)
}

// Write keeps a copy of b and passes it on.
func (a *answerRecorder) Write(b []byte) (int, error) {
	if a.status == 0 {
		a.status = http.StatusOK
	}
	a.body.Write(b)
	return a.ResponseWriter.Write(b)
}

// answer returns what was written: the status (200 if none was), the Content-Type and Location headers, and the body.
func (a *answerRecorder) answer() (int, http.Header, []byte) {
	header := http.Header{}
	for _, name := range []string{"Content-Type", "Location"} {
		if v := a.Header().Values(name); len(v) > 0 {
			header[name] = slices.Clone(v)
		}
	}
	return cmp.Or(a.status, http.StatusOK), header, a.body.Bytes()
}

// replay writes a kept answer again, marked as a replay.
func replay(w http.ResponseWriter, k *kept) {
	maps.Copy(w.Header(), k.Header)
	w.Header().Set("Idempotent-Replayed", "true")
	w.WriteHeader(k.Status)
	w.Write(k.Body)
}

// validKey reports whether key is 1 to 64 letters, digits, - or _.
func validKey(key string) bool {
	if len(key) < 1 || len(key) > 64 {
		return false
	}
	for _, r := range key {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

// idempotent answers a request that repeats an earlier one's Idempotency-Key with that request's answer.
func idempotent(keys *keyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			if !validKey(key) {
				writeError(w, r, &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: "idempotency key must be 1 to 64 letters, digits, - or _"})
				return
			}
			body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 65536))
			if _, tooLarge := errors.AsType[*http.MaxBytesError](err); tooLarge {
				writeError(w, r, &APIError{Status: http.StatusRequestEntityTooLarge, Code: "too_large", Message: "request body over 65536 bytes"})
				return
			}
			if err != nil {
				writeError(w, r, &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: "request body: " + err.Error()})
				return
			}
			c, _ := customerFrom(r.Context())
			id := keyID{Customer: c.ID, Key: key}
			k, err := keys.claim(id, body)
			switch {
			case err != nil:
				writeError(w, r, err)
			case k != nil:
				replay(w, k)
			default:
				r.Body = io.NopCloser(bytes.NewReader(body))
				rec := &answerRecorder{ResponseWriter: w}
				next.ServeHTTP(rec, r)
				status, header, out := rec.answer()
				keys.finish(id, status, header, out)
			}
		})
	}
}

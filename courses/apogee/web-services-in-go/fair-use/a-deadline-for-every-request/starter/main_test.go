package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// stalledStore is the store stuck behind Brinewatch's flood: its Get waits until the context
// it is given ends, and gives up by itself only after half a second.
type stalledStore struct {
	Store
}

// Get waits for ctx to end and returns why it ended.
func (stalledStore) Get(ctx context.Context, id string) (Booking, error) {
	select {
	case <-ctx.Done():
		return Booking{}, ctx.Err()
	case <-time.After(500 * time.Millisecond):
		return Booking{}, errors.New("the store waited half a second: the request's context had no deadline")
	}
}

// getCirrinesBooking answers like GET /bookings/bk-4521: it asks store, with the request's context, and writes what it gets.
func getCirrinesBooking(store Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := store.Get(r.Context(), "bk-4521")
		if err != nil {
			writeError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, b)
	})
}

// answer writes down a recorder's status and its ErrorBody, or its body when it holds none.
func answer(rec *httptest.ResponseRecorder) string {
	var body ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Code == "" {
		return fmt.Sprintf("%d %q", rec.Code, rec.Body)
	}
	return fmt.Sprintf("%d %s %q", rec.Code, body.Code, body.Message)
}

// serve sends GET /bookings/bk-4521 through h and returns what it answered.
func serve(h http.Handler) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/bookings/bk-4521", nil))
	return rec
}

func TestDeadlineEndsAStalledRequest(t *testing.T) {
	rec := serve(withDeadline(50 * time.Millisecond)(getCirrinesBooking(stalledStore{})))
	if got, want := answer(rec), `503 unavailable "request took too long"`; got != want {
		t.Errorf("GET /bookings/bk-4521 with a stalled store, at a 50ms deadline = %s, want %s", got, want)
	}
}

func TestDeadlineReachesTheHandler(t *testing.T) {
	var ctx context.Context
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
		w.Write([]byte("ok\n"))
	})
	rec := serve(withDeadline(10 * time.Second)(h))
	if rec.Code != http.StatusOK || rec.Body.String() != "ok\n" {
		t.Errorf("a fast handler behind a 10s deadline = %s, want 200 %q", answer(rec), "ok\n")
	}
	if ctx == nil {
		t.Fatal("the handler behind withDeadline never ran")
	}
	deadline, ok := ctx.Deadline()
	left := time.Until(deadline)
	switch {
	case !ok:
		t.Errorf("the handler's context, at a 10s deadline: Deadline() = none, want one about 10s away")
	case left > 10*time.Second || left < 9*time.Second:
		t.Errorf("the handler's context, at a 10s deadline: Deadline() = %v away, want about 10s away", left.Round(time.Millisecond))
	}
	if err := ctx.Err(); !errors.Is(err, context.Canceled) {
		t.Errorf("the handler's context once the handler returned: Err() = %v, want context.Canceled: cancel it when next returns", err)
	}
}

func TestNoDeadlineAtZero(t *testing.T) {
	ran := false
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ran = true
		if _, ok := r.Context().Deadline(); ok {
			t.Errorf("withDeadline(0): the handler's context has a deadline, want none")
		}
	})
	serve(withDeadline(0)(h))
	if !ran {
		t.Errorf("withDeadline(0): the handler never ran")
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"a deadline that passed", context.DeadlineExceeded, `503 unavailable "request took too long"`},
		{"a deadline the store wrapped", fmt.Errorf("store: %w", context.DeadlineExceeded), `503 unavailable "request took too long"`},
		{"a request canceled", context.Canceled, `503 unavailable "request canceled"`},
		{"a booking that isn't there", fmt.Errorf("booking bk-5499: %w", ErrNotFound), `404 not_found "booking bk-5499: not found"`},
		{"an APIError", &APIError{Status: 429, Code: "rate_limited", Message: "slow down: 10 requests at once, then 1 a second"}, `429 rate_limited "slow down: 10 requests at once, then 1 a second"`},
		{"anything else", errors.New("disk full"), `500 internal "internal error"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeError(rec, httptest.NewRequest(http.MethodGet, "/bookings", nil), tt.err)
			if got := answer(rec); got != tt.want {
				t.Errorf("writeError(%v) = %s, want %s", tt.err, got, tt.want)
			}
		})
	}
}

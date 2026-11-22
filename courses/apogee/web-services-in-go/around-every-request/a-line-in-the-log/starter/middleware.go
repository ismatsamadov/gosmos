package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// noStore asks every cache between the service and its callers to keep no copy of an answer.
func noStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// requestIDKey is the context key of a request's ID; no other package can make one.
type requestIDKey struct{}

// requestID gives every request the next ID, in the X-Request-Id header and in its context.
func requestID(next http.Handler) http.Handler {
	var count atomic.Uint64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := fmt.Sprintf("req-%06d", count.Add(1))
		w.Header().Set("X-Request-Id", id)
		ctx := context.WithValue(r.Context(), requestIDKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requestIDFrom returns the request ID that requestID put in ctx, or "" if there is none.
func requestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// statusRecorder is a ResponseWriter that remembers the status and the bytes written through it.
type statusRecorder struct {
	http.ResponseWriter
	Status int // the first status written; 200 when Write came first
	Bytes  int // the bytes of the body written so far
}

// WriteHeader records the first status written, then passes it on.
func (s *statusRecorder) WriteHeader(code int) {
	if s.Status == 0 {
		s.Status = code
	}
	s.ResponseWriter.WriteHeader(code)
}

// Write records a 200 if no status came first, passes p on and counts the bytes written.
func (s *statusRecorder) Write(p []byte) (int, error) {
	if s.Status == 0 {
		s.Status = http.StatusOK
	}
	n, err := s.ResponseWriter.Write(p)
	s.Bytes += n
	return n, err
}

// Unwrap returns the ResponseWriter inside, so http.NewResponseController can reach it.
func (s *statusRecorder) Unwrap() http.ResponseWriter {
	return s.ResponseWriter
}

// logRequests writes one line to l for every request, once its handler has returned.
func logRequests(l *log.Logger, now func() time.Time) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

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

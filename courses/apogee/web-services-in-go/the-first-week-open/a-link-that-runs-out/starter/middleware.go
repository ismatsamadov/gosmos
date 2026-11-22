package main

import "context"

// requestIDKey is the context key of a request's ID; no other package can make one.
type requestIDKey struct{}

// requestIDFrom returns the request ID that requestID put in ctx, or "" if there is none.
func requestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

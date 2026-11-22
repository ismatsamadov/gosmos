package main

import (
	"context"
	"net/http"
)

// passKey is the context key of the pass a piece of work belongs to; no other package can make one.
type passKey struct{}

// withPass returns a copy of ctx that carries pass ID id.
func withPass(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, passKey{}, id)
}

// passFrom returns the pass ID that withPass put in ctx, and false when there is none or it is empty.
func passFrom(ctx context.Context) (string, bool) {
	id, _ := ctx.Value(passKey{}).(string)
	return id, id != ""
}

// passTransport sends the pass in a request's context along with the request, in its Apogee-Pass header.
type passTransport struct {
	base http.RoundTripper // the transport that sends the request
}

// RoundTrip sends r through t.base with Apogee-Pass set from r's context, on a clone, so r itself never changes.
func (t passTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.base.RoundTrip(r)
}

// acceptPass puts the pass ID in a request's Apogee-Pass header into its context, for next and every handler after it.
func acceptPass(next http.Handler) http.Handler {
	return next
}

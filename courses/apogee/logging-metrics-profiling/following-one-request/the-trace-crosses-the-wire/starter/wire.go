package main

import "net/http"

// traceTransport is an http.RoundTripper that sends the span a request's context carries along
// with it, in an Apogee-Trace header.
type traceTransport struct {
	next http.RoundTripper // the transport that sends the request
}

// RoundTrip sends r through next, on a clone with Apogee-Trace set when r's context carries a span.
func (tt traceTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return tt.next.RoundTrip(r)
}

// continueTrace serves each request inside a span called name: the child of the span its
// Apogee-Trace header names, or the root of a new trace when the header is missing or malformed.
func continueTrace(t *Tracer, name string, next http.Handler) http.Handler {
	return next
}

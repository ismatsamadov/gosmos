package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// tracer is a middleware that writes its name in trace on the way in and again on the way out.
func tracer(name string, trace *[]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*trace = append(*trace, name)
			next.ServeHTTP(w, r)
			*trace = append(*trace, name)
		})
	}
}

// endpoint is the handler in the middle of a chain.
type endpoint struct{ trace *[]string }

// ServeHTTP writes "handler" in the trace and answers 204.
func (e *endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*e.trace = append(*e.trace, "handler")
	w.WriteHeader(http.StatusNoContent)
}

func TestChain(t *testing.T) {
	tests := []struct {
		name  string
		names []string
		want  string
	}{
		{"three middlewares", []string{"a", "b", "c"}, "a b c handler c b a"},
		{"one middleware", []string{"a"}, "a handler a"},
		{"none at all", nil, "handler"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var trace []string
			var mws []func(http.Handler) http.Handler
			for _, name := range tt.names {
				mws = append(mws, tracer(name, &trace))
			}
			chain(&endpoint{trace: &trace}, mws...).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/health", nil))
			call := strings.Join(append([]string{"h"}, tt.names...), ", ")
			if got := strings.Join(trace, " "); got != tt.want {
				t.Errorf("chain(%s) traced %q, want %q", call, got, tt.want)
			}
		})
	}
}

func TestChainOfNothing(t *testing.T) {
	h := &endpoint{trace: new([]string)}
	if got := chain(h); got != http.Handler(h) {
		t.Errorf("chain(h) = %v, want h itself", got)
	}
}

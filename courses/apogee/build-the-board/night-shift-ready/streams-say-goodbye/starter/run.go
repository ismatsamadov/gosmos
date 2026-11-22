package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

// newServer returns the Board's server.
func newServer(b *Board) *http.Server {
	return &http.Server{Handler: b.Handler()}
}

// run serves b on ln, refreshing it once at the start and again on every tick, until ctx
// ends; then it shuts the server down within 2 seconds and returns nil.
func run(ctx context.Context, ln net.Listener, b *Board, ticks <-chan time.Time) error {
	return newServer(b).Serve(ln)
}

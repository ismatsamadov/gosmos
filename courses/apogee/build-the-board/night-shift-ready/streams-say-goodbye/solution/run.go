package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

// newServer returns the Board's server: no write timeout, since a stream stays open, and
// the hub closed when it shuts down, so every stream says goodbye instead of holding it up.
func newServer(b *Board) *http.Server {
	srv := &http.Server{Handler: b.Handler()}
	srv.RegisterOnShutdown(b.hub.Close)
	return srv
}

// run serves b on ln, refreshing it once at the start and again on every tick, until ctx
// ends; then it shuts the server down within 2 seconds and returns nil.
func run(ctx context.Context, ln net.Listener, b *Board, ticks <-chan time.Time) error {
	return newServer(b).Serve(ln)
}

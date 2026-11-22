package main

import (
	"context"
	"errors"
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
	srv := newServer(b)
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	b.Refresh(ctx) // a refresh that fails keeps the last view
	for {
		select {
		case <-ticks:
			b.Refresh(ctx)
		case err := <-served:
			return err
		case <-ctx.Done():
			stop, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := srv.Shutdown(stop); err != nil {
				return err
			}
			if err := <-served; !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		}
	}
}

package main

import (
	"context"
	"net"
)

// run serves app on ln until ctx ends, then shuts the server down within t.Shutdown.
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error {
	srv := newServer(app.Handler, t)
	served := make(chan error, 1)
	go func() { served <- serve(srv, ln) }()
	select {
	case err := <-served:
		return err
	case <-ctx.Done():
		return stop(srv, t.Shutdown)
	}
}

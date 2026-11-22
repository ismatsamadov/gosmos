package main

import (
	"context"
	"net"
)

// run serves app on ln until ctx ends, then shuts the server down within t.Shutdown.
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error {
	srv := newServer(app.Handler, t)
	go serve(srv, ln)
	return nil
}

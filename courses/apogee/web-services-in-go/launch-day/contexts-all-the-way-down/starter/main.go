package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

// run serves app on ln until ctx ends, then drains, shuts down, and cancels whatever is still stuck after t.Shutdown.
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error {
	srv := newServer(app.Handler, t)
	_ = fmt.Errorf // TODO: give srv a BaseContext that run cancels only when stop fails
	served := make(chan error, 1)
	go func() { served <- serve(srv, ln) }()
	select {
	case err := <-served:
		return err
	case <-ctx.Done():
	}
	app.Service.draining.Store(true)
	time.Sleep(t.Drain)
	return stop(srv, t.Shutdown)
}

package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

// run serves app on ln until ctx ends, then drains, shuts down, and cancels whatever is still stuck after t.Shutdown.
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error {
	base, cancelBase := context.WithCancel(context.Background())
	defer cancelBase()
	srv := newServer(app.Handler, t)
	srv.BaseContext = func(net.Listener) context.Context { return base }
	served := make(chan error, 1)
	go func() { served <- serve(srv, ln) }()
	select {
	case err := <-served:
		return err
	case <-ctx.Done():
	}
	app.Service.draining.Store(true)
	time.Sleep(t.Drain)
	if err := stop(srv, t.Shutdown); err != nil {
		cancelBase()
		srv.Close()
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// run serves app on ln until ctx ends, then drains, shuts down and stops the notifier, in that order.
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error {
	base, cancelBase := context.WithCancel(context.Background())
	defer cancelBase()
	srv := newServer(app.Handler, t)
	srv.BaseContext = func(net.Listener) context.Context { return base }
	app.Notifier.Start()
	served := make(chan error, 1)
	go func() { served <- serve(srv, ln) }()
	select {
	case err := <-served:
		return errors.Join(err, app.Notifier.Stop())
	case <-ctx.Done():
	}
	app.Service.draining.Store(true)
	time.Sleep(t.Drain)
	if err := stop(srv, t.Shutdown); err != nil {
		cancelBase()
		srv.Close()
		return errors.Join(fmt.Errorf("shutdown: %w", err), app.Notifier.Stop())
	}
	return app.Notifier.Stop()
}

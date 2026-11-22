package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

// run serves app on ln until ctx ends, then drains for t.Drain and shuts the server down within t.Shutdown.
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error {
	srv := newServer(app.Handler, t)
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

// health answers GET /health: ok, or 503 draining while the service shuts down.
func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	if s.draining.Load() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "draining"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

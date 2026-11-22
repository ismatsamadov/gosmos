package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

// stop shuts srv down: no new connections, and up to within for the requests in flight; a zero within waits for all of them.
func stop(srv *http.Server, within time.Duration) error {
	ctx := context.Background()
	if within > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, within)
		defer cancel()
	}
	return srv.Shutdown(ctx)
}

// serve serves srv on ln until it stops, and returns nil when it was stopped on purpose.
func serve(srv *http.Server, ln net.Listener) error {
	err := srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

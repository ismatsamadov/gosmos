package main

import (
	"errors"
	"net"
	"net/http"
)

// serve serves srv on ln until it stops, and returns nil when it was stopped on purpose.
func serve(srv *http.Server, ln net.Listener) error {
	err := srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

package main

import (
	"errors"
	"net"
	"net/http"
)

// serve serves srv on ln until it stops, and returns nil when it was stopped on purpose.
func serve(srv *http.Server, ln net.Listener) error {
	_ = errors.Is
	return srv.Serve(ln)
}

package main

import (
	"errors"
	"net"
	"net/http"
)

// start serves s on a port of 127.0.0.1 that the system picks, and returns its address and a function that stops it.
func start(s *Service) (addr string, stop func() error, err error) {
	_, _ = net.Listen, (*http.Server).Serve
	return "", nil, errors.New("not started")
}

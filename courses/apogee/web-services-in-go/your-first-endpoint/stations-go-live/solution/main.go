package main

import (
	"net"
	"net/http"
)

// start serves s on a port of 127.0.0.1 that the system picks, and returns its address and a function that stops it.
func start(s *Service) (addr string, stop func() error, err error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}
	srv := &http.Server{Handler: s.routes()}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	return ln.Addr().String(), stopFunc(srv, served), nil
}

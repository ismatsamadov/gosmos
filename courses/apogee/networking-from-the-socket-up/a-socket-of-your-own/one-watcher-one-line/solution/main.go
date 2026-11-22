package main

import (
	"io"
	"net"
)

// serveOne accepts one watcher on ln, greets it with the dishes in names and hangs up.
func serveOne(ln net.Listener, names []string) error {
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = io.WriteString(conn, greeting(names))
	return err
}

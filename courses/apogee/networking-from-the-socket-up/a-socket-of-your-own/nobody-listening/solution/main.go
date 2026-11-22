package main

import (
	"errors"
	"net"
	"syscall"
)

// reach dials name's controller at addr and says in one line whether it is listening.
func reach(name, addr string) string {
	conn, err := net.Dial("tcp", addr)
	if errors.Is(err, syscall.ECONNREFUSED) {
		return name + ": not listening"
	}
	if err != nil {
		return name + ": " + err.Error()
	}
	conn.Close()
	return name + ": up"
}

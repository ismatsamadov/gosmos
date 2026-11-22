package main

import (
	"errors"
	"net"
	"syscall"
)

// reach dials name's controller at addr and says in one line whether it is listening.
func reach(name, addr string) string {
	_ = errors.Is
	_ = syscall.ECONNREFUSED
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// TODO: say "not listening" when the dial was refused
		return name + ": " + err.Error()
	}
	conn.Close()
	return name + ": up"
}

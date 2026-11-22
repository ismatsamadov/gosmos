package main

import (
	"net"
	"testing"
)

// pipePair returns two connected ends that live in memory, both closed when the test ends.
func pipePair(t testing.TB) (net.Conn, net.Conn) {
	a, b := net.Pipe()
	t.Cleanup(func() {
		a.Close()
		b.Close()
	})
	return a, b
}

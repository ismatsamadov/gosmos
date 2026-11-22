package main

import (
	"net"
	"testing"
	"time"
)

func TestListenLocal(t *testing.T) {
	first, err := listenLocal()
	if err != nil {
		t.Fatalf("listenLocal() error = %v, want a listener", err)
	}
	defer first.Close()
	checkLocal(t, "listenLocal()", first)

	second, err := listenLocal()
	if err != nil {
		t.Fatalf("a second listenLocal() while the first listens: error = %v, want a listener on another port", err)
	}
	defer second.Close()
	checkLocal(t, "a second listenLocal()", second)
	if first.Addr().String() == second.Addr().String() {
		t.Errorf("two listenLocal() calls both listen on %s, want two ports", first.Addr())
	}
}

// checkLocal fails the test unless ln listens on 127.0.0.1, on a port the system picked,
// and a client can connect to it.
func checkLocal(t *testing.T, call string, ln net.Listener) {
	t.Helper()
	addr := ln.Addr().String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("%s listens on %q: %v", call, addr, err)
	}
	if host != "127.0.0.1" {
		t.Errorf("%s listens on host %q, want 127.0.0.1, which only this Mac can reach", call, host)
	}
	if port == "7300" || port == "0" {
		t.Errorf("%s listens on port %s, want a port the system picks", call, port)
	}
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Errorf("connecting to %s: %v, want a connection", addr, err)
		return
	}
	conn.Close()
}

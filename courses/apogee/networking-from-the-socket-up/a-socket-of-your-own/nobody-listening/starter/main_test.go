package main

import (
	"net"
	"testing"
	"time"
)

func TestReach(t *testing.T) {
	sim := NewSim("YST-2", 2)
	tests := []struct {
		name, addr, want string
	}{
		{"the simulator answering as YST-2", startController(t, sim), "YST-2: up"},
		{"nobody listening", closedAddress(t), "YST-2: not listening"},
		{"an address without a port", "127.0.0.1", "YST-2: dial tcp: address 127.0.0.1: missing port in address"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reach("YST-2", tt.addr); got != tt.want {
				t.Errorf("reach(%q, %q) = %q, want %q", "YST-2", tt.addr, got, tt.want)
			}
		})
	}
	if got := closedByClient(sim, 1, time.Now().Add(time.Second)); got != 1 {
		t.Errorf("after reach found YST-2 up, the simulator saw %d connections closed by dishd, want 1: close the connection", got)
	}
}

// startController serves sim on a loopback port the system picks until the test ends, and
// returns its address.
func startController(t *testing.T, sim *Sim) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	served := make(chan error, 1)
	go func() { served <- sim.Serve(ln) }()
	t.Cleanup(func() {
		ln.Close()
		select {
		case <-served:
		case <-time.After(time.Second):
			t.Error("the simulator did not stop within a second of its listener closing")
		}
	})
	return ln.Addr().String()
}

// closedAddress returns the address of a loopback listener that has just closed, where
// nothing listens any more.
func closedAddress(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()
	return addr
}

// closedByClient waits until sim has counted want connections closed by the client, or
// until the deadline, and returns the last count it saw.
func closedByClient(sim *Sim, want int, deadline time.Time) int {
	for {
		got := sim.ClosedByClient()
		if got >= want || time.Now().After(deadline) {
			return got
		}
		time.Sleep(5 * time.Millisecond)
	}
}

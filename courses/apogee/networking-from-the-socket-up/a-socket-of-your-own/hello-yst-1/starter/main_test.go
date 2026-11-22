package main

import (
	"net"
	"testing"
	"time"
)

func TestFirstContact(t *testing.T) {
	deadline := time.Now().Add(time.Second) // every wait in this test ends by then
	tests := []struct {
		name       string
		controller string
		dish       uint8
		lastSeq    uint32
		want       string
	}{
		{"YST-1 on the bench", "YST-1", 1, 0, "YST-1 firmware 2.4, last command 0"},
		{"YST-2 after command 57", "YST-2", 2, 57, "YST-2 firmware 2.4, last command 57"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := NewSim(tt.controller, tt.dish)
			sim.LastSeq = tt.lastSeq
			sim.IdleLimit = 500 * time.Millisecond // a controller left waiting hangs up
			got, err := callFirstContact(t, startController(t, sim))
			if got != tt.want || err != nil {
				t.Fatalf("firstContact(simulator as %s) = %q, %v, want %q, nil", tt.controller, got, err, tt.want)
			}
			if n := closedByClient(sim, 1, deadline); n != 1 {
				t.Errorf("after firstContact, the simulator saw %d connections closed by dishd, want 1: close the connection when firstContact returns", n)
			}
		})
	}
	t.Run("nobody listening", func(t *testing.T) {
		addr := closedAddress(t)
		if got, err := callFirstContact(t, addr); got != "" || err == nil {
			t.Errorf("firstContact(%q) = %q, %v, want \"\" and the dial's error", addr, got, err)
		}
	})
}

// callFirstContact calls firstContact, turning a panic into a failure that says where it
// came from.
func callFirstContact(t *testing.T, addr string) (string, error) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("firstContact(%q) panicked: %v; when the dial fails conn is nil, so check the error before deferring Close", addr, p)
		}
	}()
	return firstContact(addr)
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

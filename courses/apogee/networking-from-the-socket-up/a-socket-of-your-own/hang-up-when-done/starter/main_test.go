package main

import (
	"net"
	"testing"
	"time"
)

func TestBannerOf(t *testing.T) {
	t.Run("three checks on the simulator", func(t *testing.T) {
		sim := NewSim("YST-1", 1)
		addr := startController(t, sim)
		for call := 1; call <= 3; call++ {
			got, err := callBannerOf(t, addr)
			if want := "YST-1 READY PROTO 1"; got != want || err != nil {
				t.Fatalf("bannerOf(simulator), call %d = %q, %v, want %q, nil", call, got, err, want)
			}
		}
		if got := closedByClient(sim, 3, time.Now().Add(time.Second)); got != 3 {
			t.Errorf("after 3 calls of bannerOf, the simulator saw %d connections closed by dishd, want 3: close each connection before bannerOf returns", got)
		}
	})
	t.Run("nobody listening", func(t *testing.T) {
		addr := closedAddress(t)
		if got, err := callBannerOf(t, addr); got != "" || err == nil {
			t.Errorf("bannerOf(%q) = %q, %v, want \"\" and the dial's error", addr, got, err)
		}
	})
}

// callBannerOf calls bannerOf, turning a panic into a failure that says where it came from.
func callBannerOf(t *testing.T, addr string) (string, error) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("bannerOf(%q) panicked: %v; when the dial fails conn is nil, so check the error before deferring Close", addr, p)
		}
	}()
	return bannerOf(addr)
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

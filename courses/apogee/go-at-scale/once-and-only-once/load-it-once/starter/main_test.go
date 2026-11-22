package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// errDown is what a load says while Port Ellery's file server can't be reached.
var errDown = errors.New("file server: no route to host")

// server stands in for Port Ellery's file server, where loadNetwork reads the network.
var server struct {
	loads atomic.Int64  // how many times loadNetwork has run
	asked chan struct{} // closed once every caller has asked
	down  bool          // the file server can't be reached
}

// loadNetwork stands in for reading the network from the file server: slow, since a load ends only once every caller has asked, and counted.
func loadNetwork() (*Network, error) {
	server.loads.Add(1)
	select {
	case <-server.asked:
	case <-time.After(time.Second): // never wait for a caller who doesn't come
	}
	if server.down {
		return nil, errDown
	}
	return newNetwork(), nil
}

// ask has n screens ask tb for the network at once, and returns what each got.
func ask(tb *Tables, n int) ([]*Network, []error) {
	server.loads.Store(0)
	server.asked = make(chan struct{})
	var asking atomic.Int64
	nets, errs := make([]*Network, n), make([]error, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			if asking.Add(1) == int64(n) {
				close(server.asked)
			}
			nets[i], errs[i] = tb.Network()
		})
	}
	wg.Wait()
	return nets, errs
}

func TestNetwork(t *testing.T) {
	t.Run("a hundred screens ask at 08:00", func(t *testing.T) {
		var tb Tables
		nets, errs := ask(&tb, 100)
		if n := server.loads.Load(); n != 1 {
			t.Errorf("loadNetwork ran %d times for 100 screens, want 1", n)
		}
		for i, net := range nets {
			if errs[i] != nil {
				t.Fatalf("Network() for screen %d returned the error %v, want <nil>", i, errs[i])
			}
			if net == nil || net != nets[0] {
				t.Fatalf("Network() gave screen %d a network of its own, want the one every screen shares", i)
			}
		}
		if net, err := tb.Network(); net != nets[0] || err != nil || server.loads.Load() != 1 {
			t.Errorf("a later Network() loaded again or gave another network, want the first one and no new load")
		}
	})
	t.Run("the file server is down", func(t *testing.T) {
		server.down = true
		defer func() { server.down = false }()
		var tb Tables
		nets, errs := ask(&tb, 100)
		if n := server.loads.Load(); n != 1 {
			t.Errorf("loadNetwork ran %d times for 100 screens, want 1", n)
		}
		for i := range errs {
			if !errors.Is(errs[i], errDown) || nets[i] != nil {
				t.Fatalf("Network() for screen %d returned the error %v, want no network and %q", i, errs[i], errDown)
			}
		}
	})
	t.Run("a zero Tables, one screen", func(t *testing.T) {
		var tb Tables
		nets, errs := ask(&tb, 1)
		if errs[0] != nil || nets[0] == nil || len(nets[0].Sites) != 20 {
			t.Errorf("Network() on a zero Tables returned the error %v, want 20 stations and <nil>", errs[0])
		}
	})
}

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

// fileServer stands in for Port Ellery's file server, where the network is read from.
type fileServer struct {
	loads atomic.Int64  // how many times load has run
	asked chan struct{} // closed once every caller has asked
	down  bool          // the file server can't be reached
}

// newServer returns a file server that is up, or down.
func newServer(down bool) *fileServer {
	s := &fileServer{asked: make(chan struct{}), down: down}
	close(s.asked)
	return s
}

// load stands in for reading the network from the file server: slow, since a load ends only once every caller has asked, and counted.
func (s *fileServer) load() (*Network, error) {
	s.loads.Add(1)
	select {
	case <-s.asked:
	case <-time.After(time.Second): // never wait for a caller who doesn't come
	}
	if s.down {
		return nil, errDown
	}
	return newNetwork(), nil
}

// crowd has n screens call get at once, while s is slow until all have asked, and returns what each got.
func crowd(s *fileServer, get func() (*Network, error), n int) ([]*Network, []error) {
	s.asked = make(chan struct{})
	var asking atomic.Int64
	nets, errs := make([]*Network, n), make([]error, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			if asking.Add(1) == int64(n) {
				close(s.asked)
			}
			nets[i], errs[i] = get()
		})
	}
	wg.Wait()
	return nets, errs
}

func TestNetworkOnce(t *testing.T) {
	t.Run("a hundred screens ask at 08:00", func(t *testing.T) {
		s := newServer(false)
		nets, errs := crowd(s, networkOnce(s.load), 100)
		if n := s.loads.Load(); n != 1 {
			t.Errorf("load ran %d times for 100 screens, want 1", n)
		}
		for i, net := range nets {
			if errs[i] != nil {
				t.Fatalf("screen %d got the error %v, want <nil>", i, errs[i])
			}
			if net == nil || net != nets[0] {
				t.Fatalf("screen %d got a network of its own, want the one every screen shares", i)
			}
		}
	})
	t.Run("the file server is down", func(t *testing.T) {
		s := newServer(true)
		get := networkOnce(s.load)
		_, errs := crowd(s, get, 100)
		for i, err := range errs {
			if !errors.Is(err, errDown) {
				t.Fatalf("screen %d got the error %v, want %q", i, err, errDown)
			}
		}
		if net, err := get(); net != nil || !errors.Is(err, errDown) {
			t.Errorf("a call after the failed load got the error %v, want %q again", err, errDown)
		}
		if n := s.loads.Load(); n != 1 {
			t.Errorf("load ran %d times for 101 calls while the server was down, want 1: no second try", n)
		}
	})
	t.Run("two tables, each loaded on its own", func(t *testing.T) {
		a, b := newServer(false), newServer(false)
		getA, getB := networkOnce(a.load), networkOnce(b.load)
		netA, _ := getA()
		netB, _ := getB()
		againA, _ := getA()
		againB, _ := getB()
		if na, nb := a.loads.Load(), b.loads.Load(); na != 1 || nb != 1 {
			t.Errorf("loads after two calls each = %d and %d, want 1 and 1", na, nb)
		}
		if againA != netA || againB != netB {
			t.Errorf("a second call of a function from networkOnce gave another network, want the first one again")
		}
		if netA == nil || netA == netB {
			t.Errorf("networkOnce(a.load) and networkOnce(b.load) gave the same network, want one each")
		}
	})
}

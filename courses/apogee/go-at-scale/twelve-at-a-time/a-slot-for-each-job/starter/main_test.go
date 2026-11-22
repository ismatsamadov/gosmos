package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// guard is how long a test waits for what a right stage does at once: one second in all.
type guard struct{ end time.Time }

// newGuard starts a guard's second.
func newGuard() guard { return guard{end: time.Now().Add(time.Second)} }

// wait is ready when the guard's second is over, and never sooner than 50 ms from now.
func (g guard) wait() <-chan time.Time {
	return time.After(max(time.Until(g.end), 50*time.Millisecond))
}

// meter counts the jobs in flight and keeps the most there have been at once.
type meter struct{ now, most atomic.Int64 }

// start counts a job in.
func (m *meter) start() {
	n := m.now.Add(1)
	for top := m.most.Load(); n > top && !m.most.CompareAndSwap(top, n); top = m.most.Load() {
	}
}

// stop counts a job out.
func (m *meter) stop() { m.now.Add(-1) }

// fleet returns the eighty satellites' elements in fleet order, with only their IDs filled in.
func fleet() []Elements {
	ids := []string{"pipit-1", "pipit-2", "pipit-3", "furrow-a", "furrow-b", "brine-3", "cirrine-1", "thistle-2"}
	for n := 2; n <= 13; n++ {
		ids = append(ids, fmt.Sprintf("cirrine-%d", n))
	}
	for n := 4; n <= 39; n++ {
		ids = append(ids, fmt.Sprintf("brine-%d", n))
	}
	for c := 'c'; c <= 'z'; c++ {
		ids = append(ids, fmt.Sprintf("furrow-%c", c))
	}
	els := make([]Elements, len(ids))
	for i, id := range ids {
		els[i] = Elements{SatID: id}
	}
	return els
}

func TestSlots(t *testing.T) {
	g := newGuard()
	t.Run("twelve planner slots taken", func(t *testing.T) {
		s := newSlots(12)
		for i := range 12 {
			if err := s.Acquire(context.Background()); err != nil {
				t.Fatalf("Acquire %d of 12 = %v, want nil", i+1, err)
			}
		}
		if len(s) != 12 {
			t.Errorf("slots taken after 12 Acquires = %d, want 12", len(s))
		}
	})
	t.Run("a slot given back", func(t *testing.T) {
		s := newSlots(2)
		s <- struct{}{}
		s <- struct{}{}
		s.Release()
		if len(s) != 1 {
			t.Errorf("slots taken after a Release = %d of 2, want 1", len(s))
		}
	})
	t.Run("the thirteenth job takes the slot given back", func(t *testing.T) {
		s := newSlots(12)
		for range 12 {
			s <- struct{}{}
		}
		got := make(chan error, 1)
		go func() { got <- s.Acquire(context.Background()) }()
		s.Release()
		select {
		case err := <-got:
			if err != nil {
				t.Errorf("Acquire after a Release = %v, want nil", err)
			}
		case <-g.wait():
			t.Fatalf("Acquire is still waiting after a Release gave a slot back, want it to take that slot")
		}
		if len(s) != 12 {
			t.Errorf("slots taken = %d, want 12: the thirteenth job holds the slot given back", len(s))
		}
	})
	t.Run("a canceled run waits for no slot", func(t *testing.T) {
		s := newSlots(2)
		s <- struct{}{}
		s <- struct{}{}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		got := make(chan error, 1)
		go func() { got <- s.Acquire(ctx) }()
		select {
		case err := <-got:
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Acquire(canceled ctx) on full slots = %v, want context.Canceled", err)
			}
		case <-g.wait():
			t.Errorf("Acquire(canceled ctx) on full slots is still waiting, want context.Canceled at once")
		}
	})
}

func TestPredictAll(t *testing.T) {
	g := newGuard()
	els := fleet()
	var m meter
	arrived := make(chan string, len(els))
	release := make(chan struct{})
	// hold is a prediction that waits until the test lets every prediction go.
	hold := func(ctx context.Context, e Elements) ([]PassView, error) {
		m.start()
		defer m.stop()
		arrived <- e.SatID
		<-release
		return nil, nil
	}
	type result struct {
		passes map[string][]PassView
		err    error
	}
	done := make(chan result, 1)
	go func() {
		passes, err := predictAll(context.Background(), els, 12, hold)
		done <- result{passes, err}
	}()
	for i := range 12 {
		select {
		case <-arrived:
		case <-g.wait():
			close(release)
			t.Fatalf("predictions in flight = %d, want 12 at once", i)
		}
	}
	// With twelve predictions held, right slots start no thirteenth: 50 ms gives wrong ones time to.
	<-time.After(50 * time.Millisecond)
	close(release)
	var r result
	select {
	case r = <-done:
	case <-g.wait():
		t.Fatalf("predictAll(80 satellites, 12 slots) did not return once every prediction was let go")
	}
	if n := m.most.Load(); n != 12 {
		t.Errorf("most in flight = %d, want 12", n)
	}
	if r.err != nil || len(r.passes) != 80 {
		t.Errorf("predictAll(80 satellites, 12 slots) = %d satellites, %v, want 80, <nil>", len(r.passes), r.err)
	}
}

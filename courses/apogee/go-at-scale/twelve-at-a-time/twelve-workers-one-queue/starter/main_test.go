package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
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

// queue returns a closed channel holding els, in order: every job of a run.
func queue(els []Elements) <-chan Elements {
	ch := make(chan Elements, len(els))
	for _, e := range els {
		ch <- e
	}
	close(ch)
	return ch
}

func TestRunPool(t *testing.T) {
	g := newGuard()
	t.Run("eighty satellites, twelve workers", func(t *testing.T) {
		before := runtime.NumGoroutine()
		var m meter
		var calls atomic.Int64
		arrived := make(chan struct{}, 80)
		release := make(chan struct{})
		// hold is a job that waits until the test lets every job go.
		hold := func(ctx context.Context, e Elements) error {
			m.start()
			defer m.stop()
			calls.Add(1)
			arrived <- struct{}{}
			<-release
			return nil
		}
		done := make(chan error, 1)
		go func() { done <- runPool(context.Background(), queue(fleet()), 12, hold) }()
		for i := range 12 {
			select {
			case <-arrived:
			case <-g.wait():
				close(release)
				t.Fatalf("jobs in flight = %d, want 12 at once", i)
			}
		}
		// With twelve jobs held, a right pool starts nothing more: 50 ms gives a wrong one time to.
		<-time.After(50 * time.Millisecond)
		running := runtime.NumGoroutine() - before - 1 // every goroutine but the one running runPool
		close(release)
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("runPool(80 jobs, 12 workers) = %v, want <nil>", err)
			}
		case <-g.wait():
			t.Fatalf("runPool(80 jobs, 12 workers) did not return once every job was let go")
		}
		if running > 12 {
			t.Errorf("goroutines while 12 jobs were held = %d, want 12: one for each worker", running)
		}
		if n := m.most.Load(); n != 12 {
			t.Errorf("most jobs in flight = %d, want 12", n)
		}
		if n := calls.Load(); n != 80 {
			t.Errorf("do called %d times for 80 jobs, want 80", n)
		}
	})
	t.Run("cirrine-9 has no elements yet", func(t *testing.T) {
		missing := errors.New("no elements for cirrine-9")
		var finished atomic.Int64
		do := func(ctx context.Context, e Elements) error {
			defer finished.Add(1)
			if e.SatID == "cirrine-9" {
				return missing
			}
			return nil
		}
		done := make(chan error, 1)
		go func() { done <- runPool(context.Background(), queue(fleet()[8:20]), 4, do) }() // the twelve new Cirrines
		select {
		case err := <-done:
			if !errors.Is(err, missing) {
				t.Errorf("runPool(12 Cirrines, 4 workers) = %v, want %q", err, missing)
			}
			if n := finished.Load(); n != 12 {
				t.Errorf("jobs finished when runPool returned = %d, want all 12", n)
			}
		case <-g.wait():
			t.Fatalf("runPool(12 Cirrines, 4 workers) did not return")
		}
	})
	t.Run("no satellites queued", func(t *testing.T) {
		done := make(chan error, 1)
		go func() {
			done <- runPool(context.Background(), queue(nil), 12, func(context.Context, Elements) error { return nil })
		}()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("runPool(no jobs) = %v, want <nil>", err)
			}
		case <-g.wait():
			t.Fatalf("runPool(no jobs) did not return, want <nil> at once")
		}
	})
}

package main

import (
	"sync"
	"testing"
)

// stop calls r.Stop and returns what it panicked with, or nil.
func stop(r *Run) (panicked any) {
	defer func() { panicked = recover() }()
	r.Stop()
	return nil
}

// closed reports whether ch is closed, without waiting.
func closed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func TestStop(t *testing.T) {
	t.Run("the Board's shutdown, then a newer plan", func(t *testing.T) {
		r := newRun()
		if p := stop(r); p != nil {
			t.Fatalf("first Stop panicked: %v", p)
		}
		if p := stop(r); p != nil {
			t.Errorf("second Stop panicked: %v", p)
		}
		if !closed(r.Done()) {
			t.Errorf("Done() is still open after Stop, want it closed")
		}
	})
	t.Run("ten callers stop the run at once", func(t *testing.T) {
		r := newRun()
		start := make(chan struct{})
		panics := make(chan any, 10)
		var wg sync.WaitGroup
		for range 10 {
			wg.Go(func() {
				<-start
				if p := stop(r); p != nil {
					panics <- p
				}
			})
		}
		close(start)
		wg.Wait()
		close(panics)
		if n := len(panics); n > 0 {
			t.Errorf("Stop from 10 goroutines at once: %d panicked (%v), want none", n, <-panics)
		}
		if !closed(r.Done()) {
			t.Errorf("Done() is still open after Stop, want it closed")
		}
	})
	t.Run("a run nobody has stopped", func(t *testing.T) {
		r := newRun()
		if closed(r.Done()) {
			t.Errorf("Done() is closed before any Stop, want it open while the run goes on")
		}
	})
}

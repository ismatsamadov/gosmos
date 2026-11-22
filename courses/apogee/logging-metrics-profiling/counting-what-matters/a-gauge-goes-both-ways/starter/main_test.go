package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// failure is the error Cirrine-1's pass ended with at 02:02:39 on 8 November.
var failure = fmt.Errorf("pass cirrine-1/LAN/16402: %w", context.DeadlineExceeded)

func TestTrack(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"a flown pass", nil},
		{"a failed pass", failure},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g Gauge
			called := false
			var inside int64
			err := track(&g, func() error {
				called = true
				inside = g.Value()
				return tt.err
			})
			if !called {
				t.Fatalf("track(g, run) returned %v without calling run", err)
			}
			if inside != 1 {
				t.Errorf("g.Value() inside run = %d, want 1", inside)
			}
			if err != tt.err {
				t.Errorf("track(g, run) = %v, want run's error, %v", err, tt.err)
			}
			if got := g.Value(); got != 0 {
				t.Errorf("g.Value() after track(g, run) = %d, want 0", got)
			}
		})
	}
}

// TestTrackPanics runs a pass whose run panics, as a hung radio's code might: only a
// deferred call still brings the gauge down while the panic goes on to whoever recovers it.
func TestTrackPanics(t *testing.T) {
	var g Gauge
	func() {
		defer func() { _ = recover() }()
		_ = track(&g, func() error { panic("radio hung") })
	}()
	if got := g.Value(); got != 0 {
		t.Errorf("g.Value() after run panicked = %d, want 0", got)
	}
}

// TestTwentyAtOnce runs 20 passes at once, every other one failing. While all 20 are inside
// run the gauge reads 20; once they have all ended it reads 0.
func TestTwentyAtOnce(t *testing.T) {
	var g Gauge
	var entered atomic.Int64
	allIn := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			_ = track(&g, func() error {
				if entered.Add(1) == 20 {
					close(allIn)
				}
				<-release
				if i%2 == 1 {
					return failure
				}
				return nil
			})
		})
	}
	returned := make(chan struct{}) // closed if every call returns, as when track never calls run
	go func() { wg.Wait(); close(returned) }()
	select {
	case <-allIn:
	case <-returned:
		t.Fatalf("20 calls to track(g, run) ran run %d times, want 20", entered.Load())
	case <-time.After(5 * time.Second): // a safety net only: 20 goroutines start within microseconds
		close(release)
		wg.Wait()
		t.Fatalf("20 calls to track(g, run) ran run %d times, want 20", entered.Load())
	}
	during := g.Value()
	close(release)
	wg.Wait()
	if during != 20 {
		t.Errorf("g.Value() with 20 passes inside run = %d, want 20", during)
	}
	if got := g.Value(); got != 0 {
		t.Errorf("g.Value() after 20 passes, 10 of them failed = %d, want 0", got)
	}
}

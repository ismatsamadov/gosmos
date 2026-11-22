package main

import (
	"context"
	"errors"
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

func TestPredictAll(t *testing.T) {
	n := newNetwork()
	t.Run("take-over day, eighty satellites over twenty stations", func(t *testing.T) {
		day := time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)
		var m meter
		// overDay predicts e's passes over every station on 28 June with the pass model.
		overDay := func(ctx context.Context, e Elements) ([]PassView, error) {
			m.start()
			defer m.stop()
			var passes []PassView
			for _, s := range n.Sites {
				passes = append(passes, predict(e, s, day, day.Add(24*time.Hour))...)
			}
			return passes, nil
		}
		type result struct {
			passes map[string][]PassView
			err    error
		}
		done := make(chan result, 1)
		go func() {
			passes, err := predictAll(context.Background(), n.Elements, 12, overDay)
			done <- result{passes, err}
		}()
		var r result
		select {
		case r = <-done:
		case <-time.After(5 * time.Second): // real predictions: far more time than a day of them takes
			t.Fatalf("predictAll(28 June, 12 slots) did not return")
		}
		if r.err != nil {
			t.Fatalf("predictAll(28 June, 12 slots) error = %v, want <nil>", r.err)
		}
		total := 0
		for _, passes := range r.passes {
			total += len(passes)
		}
		if len(r.passes) != 80 || len(r.passes["furrow-a"]) != 105 || total != 7485 {
			t.Errorf("predictAll(28 June, 12 slots) = %d satellites, Furrow-A with %d passes, %d in all, want 80, 105 and 7485", len(r.passes), len(r.passes["furrow-a"]), total)
		}
		if most := m.most.Load(); most > 12 {
			t.Errorf("most predictions at once = %d, want 12 at most", most)
		}
	})
	g := newGuard() // from here, every wait is for what a right group does at once
	t.Run("brine-17 has no elements", func(t *testing.T) {
		missing := errors.New("no elements for brine-17")
		const failing = 33 // brine-17's place in fleet order
		index := map[string]int{}
		for i, e := range n.Elements {
			index[e.SatID] = i
		}
		var held, late atomic.Int64
		allHeld := make(chan struct{})
		lateStart := make(chan struct{}, 46)
		// predict returns at once for the first 22 satellites, holds the next 11 until their context ends, fails for brine-17 once those 11 are running, and counts any prediction started after it.
		predict := func(ctx context.Context, e Elements) ([]PassView, error) {
			switch i := index[e.SatID]; {
			case i < failing-11:
				return nil, nil
			case i < failing:
				if held.Add(1) == 11 {
					close(allHeld)
				}
				<-ctx.Done()
				return nil, ctx.Err()
			case i == failing:
				select {
				case <-allHeld:
				case <-time.After(time.Second):
				}
				return nil, missing
			default:
				late.Add(1)
				lateStart <- struct{}{}
				return nil, nil
			}
		}
		parent, stop := context.WithCancel(context.Background())
		defer stop()
		done := make(chan error, 1)
		go func() {
			_, err := predictAll(parent, n.Elements, 12, predict)
			done <- err
		}()
		var err error
		returned := false
		select {
		case err = <-done:
			returned = true
		case <-lateStart:
			t.Errorf("a prediction started after brine-17 failed, want none: the run went on")
		case <-g.wait():
			t.Errorf("a second after brine-17 failed, predictAll was still waiting on predictions that never saw their context end")
		}
		if !returned {
			stop() // the test cancels the run itself, so every prediction ends
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatalf("predictAll did not return even after the run was canceled")
			}
			return
		}
		if !errors.Is(err, missing) {
			t.Errorf("predictAll(80 satellites, brine-17 failing) error = %v, want %q", err, missing)
		}
		if k := late.Load(); k != 0 {
			t.Errorf("predictions started after brine-17 failed = %d, want 0", k)
		}
	})
}

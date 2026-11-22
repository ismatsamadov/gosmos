package main

import (
	"context"
	"errors"
	"sync"
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

// held is a set of uploads that each wait until the test lets them all go.
type held struct {
	arrived chan string
	release chan struct{}
}

// newHeld makes room for n held uploads.
func newHeld(n int) *held {
	return &held{arrived: make(chan string, n), release: make(chan struct{})}
}

// upload returns an upload to station that tells the test it has started and waits to be let go, counted by m.
func (h *held) upload(station string, m ...*meter) func() error {
	return func() error {
		for _, m := range m {
			m.start()
			defer m.stop()
		}
		h.arrived <- station
		<-h.release
		return nil
	}
}

// await waits until n held uploads have started, and returns how many did before the guard ended.
func (h *held) await(g guard, n int) int {
	for i := range n {
		select {
		case <-h.arrived:
		case <-g.wait():
			return i
		}
	}
	return n
}

func TestUploads(t *testing.T) {
	g := newGuard()
	t.Run("six uploads each to Calidra Atoll and Varsk Island at once", func(t *testing.T) {
		u := newUploads(2)
		cal, vsk, together := new(meter), new(meter), new(meter)
		h := newHeld(12)
		var wg sync.WaitGroup
		errs := make(chan error, 12)
		for range 6 {
			wg.Go(func() { errs <- u.Do(context.Background(), "CAL", h.upload("CAL", cal, together)) })
			wg.Go(func() { errs <- u.Do(context.Background(), "VSK", h.upload("VSK", vsk, together)) })
		}
		n := h.await(g, 4)
		if n == 4 {
			// With two uploads held at each station, right slots start no more: 50 ms gives wrong ones time to.
			<-time.After(50 * time.Millisecond)
		}
		close(h.release)
		wg.Wait()
		close(errs)
		if n < 4 {
			t.Fatalf("uploads in flight together = %d, want 4: two to each station", n)
		}
		for err := range errs {
			if err != nil {
				t.Errorf("Do(ctx, station, an upload that succeeds) = %v, want <nil>", err)
			}
		}
		for _, s := range []struct {
			code string
			m    *meter
		}{{"CAL", cal}, {"VSK", vsk}} {
			if most := s.m.most.Load(); most != 2 {
				t.Errorf("most uploads at once to %s = %d, want 2", s.code, most)
			}
		}
		if most := together.most.Load(); most != 4 {
			t.Errorf("most uploads at once to both stations = %d, want 4", most)
		}
	})
	t.Run("a canceled upload waits for no slot", func(t *testing.T) {
		u := newUploads(2)
		h := newHeld(2)
		defer close(h.release)
		for range 2 {
			go u.Do(context.Background(), "CAL", h.upload("CAL"))
		}
		if n := h.await(g, 2); n < 2 {
			t.Fatalf("uploads to CAL in flight = %d, want 2", n)
		}
		ctx, cancel := context.WithCancel(context.Background())
		var ran atomic.Bool
		done := make(chan error, 1)
		go func() {
			done <- u.Do(ctx, "CAL", func() error {
				ran.Store(true)
				return nil
			})
		}()
		cancel()
		select {
		case err := <-done:
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Do(canceled ctx, CAL with both slots taken) = %v, want context.Canceled", err)
			}
		case <-g.wait():
			t.Fatalf("Do(canceled ctx, CAL with both slots taken) is still waiting, want context.Canceled")
		}
		if ran.Load() {
			t.Errorf("an upload ran though its run was canceled while it waited for a slot")
		}
	})
	t.Run("a failed upload gives its slot back", func(t *testing.T) {
		u := newUploads(2)
		dropped := errors.New("vsk: line dropped")
		if err := u.Do(context.Background(), "VSK", func() error { return dropped }); !errors.Is(err, dropped) {
			t.Errorf("Do(ctx, VSK, a failing upload) = %v, want %q", err, dropped)
		}
		h := newHeld(2)
		defer close(h.release)
		for range 2 {
			go u.Do(context.Background(), "VSK", h.upload("VSK"))
		}
		if n := h.await(g, 2); n < 2 {
			t.Errorf("uploads to VSK in flight after a failed one = %d, want 2: the failed upload's slot is back", n)
		}
	})
}

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

// waitFor returns what Wait returns, or reports that it is still waiting once the guard ends.
func waitFor(g guard, grp *Group) (err error, returned bool) {
	done := make(chan error, 1)
	go func() { done <- grp.Wait() }()
	select {
	case err := <-done:
		return err, true
	case <-g.wait():
		return nil, false
	}
}

func TestGroup(t *testing.T) {
	g := newGuard()
	missing := errors.New("no elements for brine-17")
	t.Run("twenty predictions, four slots", func(t *testing.T) {
		grp, _ := newGroup(context.Background(), 4)
		var m meter
		var ran atomic.Int64
		arrived := make(chan struct{}, 20)
		release := make(chan struct{})
		// hold is a job that waits until the test lets every job go.
		hold := func(ctx context.Context) error {
			m.start()
			defer m.stop()
			ran.Add(1)
			arrived <- struct{}{}
			<-release
			return nil
		}
		done := make(chan error, 1)
		go func() {
			for range 20 {
				grp.Go(hold)
			}
			done <- grp.Wait()
		}()
		for i := range 4 {
			select {
			case <-arrived:
			case <-g.wait():
				close(release)
				t.Fatalf("jobs in flight = %d, want 4 at once", i)
			}
		}
		close(release)
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("Wait() = %v after twenty jobs that succeed, want <nil>", err)
			}
		case <-g.wait():
			t.Fatalf("Wait() did not return once every job was let go")
		}
		if n := ran.Load(); n != 20 {
			t.Errorf("jobs run = %d, want all 20", n)
		}
		if n := m.most.Load(); n > 4 {
			t.Errorf("most jobs in flight = %d, want 4", n)
		}
	})
	kept := t.Run("brine-17's prediction fails", func(t *testing.T) {
		grp, ctx := newGroup(context.Background(), 4)
		grp.Go(func(context.Context) error { return missing })
		err, returned := waitFor(g, grp)
		if !returned {
			t.Fatalf("Wait() did not return after brine-17 failed")
		}
		if !errors.Is(err, missing) {
			t.Errorf("Wait() = %v, want %q: a job's error fails the group", err, missing)
		}
		if cause := context.Cause(ctx); !errors.Is(cause, missing) {
			t.Errorf("context.Cause(the jobs' context) = %v, want %q, the failure that ended it", cause, missing)
		}
	})
	if !kept {
		return // with no failure kept, the checks below could only wait out the guard
	}
	t.Run("brine-17 fails while two predictions wait on their context", func(t *testing.T) {
		parent, stop := context.WithCancel(context.Background())
		defer stop()
		grp, _ := newGroup(parent, 4)
		started := make(chan struct{}, 2)
		sawEnd := make(chan struct{}, 2)
		// hold is a job that runs until its context ends.
		hold := func(ctx context.Context) error {
			started <- struct{}{}
			<-ctx.Done()
			sawEnd <- struct{}{}
			return ctx.Err()
		}
		grp.Go(hold) // brine-15
		grp.Go(hold) // brine-16
		for i := range 2 {
			select {
			case <-started:
			case <-g.wait():
				stop()
				t.Fatalf("jobs running = %d, want 2 before brine-17's", i)
			}
		}
		grp.Go(func(context.Context) error { return missing })
		for i := 0; i < 2; i++ {
			select {
			case <-sawEnd:
				continue
			case <-g.wait():
			}
			t.Errorf("after brine-17 failed, %d of 2 running jobs saw their context end, want 2", i)
			stop() // the test cancels instead, so the held jobs end
			break
		}
		var late atomic.Bool
		grp.Go(func(context.Context) error {
			late.Store(true)
			return nil
		})
		err, returned := waitFor(g, grp)
		if !returned {
			stop()
			t.Fatalf("Wait() did not return after brine-17 failed")
		}
		if !errors.Is(err, missing) {
			t.Errorf("Wait() = %v, want %q, the first failure", err, missing)
		}
		if late.Load() {
			t.Errorf("a job given to Go after brine-17 failed ran, want it never to run")
		}
	})
	t.Run("brine-17 fails, then brine-18", func(t *testing.T) {
		grp, _ := newGroup(context.Background(), 2)
		later := errors.New("brine-18: prediction stopped")
		grp.Go(func(context.Context) error { return missing })
		grp.Go(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
			case <-time.After(time.Second): // never wait for a group that doesn't cancel
			}
			return later
		})
		err, returned := waitFor(g, grp)
		if !returned {
			t.Fatalf("Wait() did not return after two failures")
		}
		if !errors.Is(err, missing) {
			t.Errorf("Wait() = %v, want %q, the first of two failures", err, missing)
		}
	})
	t.Run("the run canceled before any prediction starts", func(t *testing.T) {
		parent, stop := context.WithCancel(context.Background())
		stop()
		grp, _ := newGroup(parent, 4)
		var ran atomic.Int64
		given := make(chan struct{})
		go func() {
			defer close(given)
			for range 6 { // more jobs than slots: each one refused must give its slot back
				grp.Go(func(context.Context) error {
					ran.Add(1)
					return nil
				})
			}
		}()
		select {
		case <-given:
		case <-g.wait():
			t.Fatalf("Go(6 jobs of a canceled run, 4 slots) is still waiting for a slot, want each refused job to give its slot back")
		}
		err, returned := waitFor(g, grp)
		if !returned {
			t.Fatalf("Wait() did not return for a canceled run")
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Wait() = %v, want context.Canceled: a job Go refuses is a failure too", err)
		}
		if n := ran.Load(); n != 0 {
			t.Errorf("jobs run after the run was canceled = %d, want 0", n)
		}
	})
}

package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"runtime"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
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

// goroutinesLeft waits until no more than before goroutines run, or the guard ends, and returns how many more run.
func goroutinesLeft(g guard, before int) int {
	end := time.Now().Add(max(time.Until(g.end), 50*time.Millisecond))
	for {
		n := runtime.NumGoroutine() - before
		if n <= 0 || time.Now().After(end) {
			return max(n, 0)
		}
		time.Sleep(time.Millisecond)
	}
}

// takeOverDay is midnight UTC on Thursday 28 June 2029, where the tests set a bubble's clock.
var takeOverDay = time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)

// clockAt returns a clock that reads day when the bubble starts and moves with its fake time.
func clockAt(day time.Time) func() time.Time {
	shift := day.Sub(time.Now())
	return func() time.Time { return time.Now().Add(shift).UTC() }
}

// rehearsal is passplan 2's settings, as Tuesday's rehearsal ran them.
var rehearsal = Config{Slots: 12, Days: 1, MinElevDeg: 10, Turnaround: 2 * time.Minute, Replan: 10 * time.Minute, RunLimit: 30 * time.Second}

// marked returns a fake run's schedule: one contact, named for the time of the event it was planned for.
func marked(at time.Time) *Schedule {
	return &Schedule{Contacts: []Contact{{PassView: PassView{ID: at.Format(time.TimeOnly)}}}}
}

// shown describes a schedule as the tests print it.
func shown(s *Schedule) string {
	if s == nil {
		return "none"
	}
	return fmt.Sprintf("run %d's, for the event at %s", s.Version, s.Contacts[0].ID)
}

// send gives Run an event, and fails the test if Run doesn't take it within a fake minute.
func send(t *testing.T, events chan<- Event, ev Event) {
	t.Helper()
	select {
	case events <- ev:
	case <-time.After(time.Minute):
		t.Fatalf("Run did not take the event at %s", ev.At.Format(time.TimeOnly))
	}
}

// over closes events and waits up to a fake hour for Run to return nil.
func over(t *testing.T, events chan Event, done <-chan error) {
	t.Helper()
	close(events)
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run = %v once events closed, want <nil>", err)
		}
	case <-time.After(time.Hour):
		t.Fatalf("Run did not return an hour after events closed")
	}
}

func TestRun(t *testing.T) {
	t.Run("a stow five seconds into a twenty-second run", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			now := clockAt(takeOverDay)
			pub := &Publisher{}
			r := newReplanner(nil, pub, rehearsal, now)
			var mu sync.Mutex
			var canceled []string // the events whose runs saw their context end
			r.plan = func(ctx context.Context, cur *Schedule, at time.Time, stows []Stow) (*Schedule, error) {
				select {
				case <-time.After(20 * time.Second):
					return marked(at), nil
				case <-ctx.Done():
					mu.Lock()
					canceled = append(canceled, at.Format(time.TimeOnly))
					mu.Unlock()
					return nil, ctx.Err()
				}
			}
			events, done := make(chan Event), make(chan error, 1)
			go func() { done <- r.Run(t.Context(), events) }()
			send(t, events, Event{At: now()}) // a tick at 00:00:00
			time.Sleep(5 * time.Second)
			send(t, events, Event{At: now(), Stow: &Stow{Station: "CWR", From: now(), Reason: "wind 74 km/h"}})
			time.Sleep(16 * time.Second)
			if cur := pub.Current(); cur != nil {
				t.Errorf("at 00:00:21, current schedule = %s, want none yet: the stow's news cancels the tick's run", shown(cur))
			}
			time.Sleep(5 * time.Second)
			if cur := pub.Current(); cur == nil || cur.Version != 2 || cur.Contacts[0].ID != "00:00:05" {
				t.Errorf("at 00:00:26, current schedule = %s, want run 2's, for the event at 00:00:05", shown(cur))
			}
			mu.Lock()
			got := slices.Clone(canceled)
			mu.Unlock()
			if !slices.Equal(got, []string{"00:00:00"}) {
				t.Errorf("runs whose context ended = %v, want [00:00:00]: news cancels the run in progress", got)
			}
			over(t, events, done)
		})
	})
	t.Run("a tick that finds nothing new", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			now := clockAt(takeOverDay)
			v7 := marked(takeOverDay.Add(-time.Hour))
			v7.Version = 7
			pub := &Publisher{}
			pub.Offer(v7)
			r := newReplanner(nil, pub, rehearsal, now)
			r.plan = func(ctx context.Context, cur *Schedule, at time.Time, stows []Stow) (*Schedule, error) {
				time.Sleep(time.Second)
				if at.Equal(takeOverDay) {
					return &Schedule{Contacts: slices.Clone(cur.Contacts)}, nil // nothing new
				}
				return marked(at), nil
			}
			events, done := make(chan Event), make(chan error, 1)
			go func() { done <- r.Run(t.Context(), events) }()
			send(t, events, Event{At: now()}) // a tick at 00:00:00
			time.Sleep(2 * time.Second)
			if cur := pub.Current(); cur != v7 {
				t.Errorf("after a run that found nothing new, current schedule = %s, want run 7's still", shown(cur))
			}
			send(t, events, Event{At: now()}) // a tick at 00:00:02
			time.Sleep(2 * time.Second)
			if cur := pub.Current(); cur == nil || cur.Version != 9 || cur.Contacts[0].ID != "00:00:02" {
				t.Errorf("at 00:00:04, current schedule = %s, want run 9's, for the event at 00:00:02", shown(cur))
			}
			over(t, events, done)
		})
	})
	t.Run("thirty ticks a second apart, their runs ending in a shuffled order", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			now := clockAt(takeOverDay)
			pub := &Publisher{}
			r := newReplanner(nil, pub, rehearsal, now)
			takes := rand.New(rand.NewPCG(1, 2)).Perm(30) // the run for second k takes takes[k-1]+1 seconds
			r.plan = func(ctx context.Context, cur *Schedule, at time.Time, stows []Stow) (*Schedule, error) {
				k := min(max(int(at.Sub(takeOverDay)/time.Second), 1), len(takes))
				time.Sleep(time.Duration(takes[k-1]+1) * time.Second) // it finishes its work, canceled or not
				return marked(at), nil
			}
			events, done := make(chan Event, 30), make(chan error, 1) // room for every event: none waits for Run
			go func() { done <- r.Run(t.Context(), events) }()
			for range 30 {
				time.Sleep(time.Second)
				send(t, events, Event{At: now()})
			}
			time.Sleep(8 * time.Minute) // every run has ended, even one run after another, before Run's tick at 00:10
			if cur := pub.Current(); cur == nil || cur.Version != 30 || cur.Contacts[0].ID != "00:00:30" {
				t.Errorf("after thirty runs ending in a shuffled order, current schedule = %s, want run 30's, for the event at 00:00:30", shown(cur))
			}
			over(t, events, done)
		})
	})
	t.Run("ten minutes with no news, then events closes", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			now := clockAt(takeOverDay)
			pub := &Publisher{}
			r := newReplanner(nil, pub, rehearsal, now)
			r.plan = func(ctx context.Context, cur *Schedule, at time.Time, stows []Stow) (*Schedule, error) {
				select {
				case <-time.After(time.Second):
					return marked(at), nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
			events, done := make(chan Event), make(chan error, 1)
			go func() { done <- r.Run(t.Context(), events) }()
			time.Sleep(10*time.Minute + 500*time.Millisecond) // the tick's run is half done
			over(t, events, done)
			if cur := pub.Current(); cur == nil || cur.Version != 1 || cur.Contacts[0].ID != "00:10:00" {
				t.Errorf("once Run has returned, current schedule = %s, want run 1's, for the tick at 00:10:00: the last run finishes before Run returns", shown(cur))
			}
		})
	})
	t.Run("the context ends with a run in flight", func(t *testing.T) {
		g := newGuard()
		before := runtime.NumGoroutine()
		pub := &Publisher{}
		r := newReplanner(nil, pub, rehearsal, time.Now)
		started := make(chan struct{}, 3)
		r.plan = func(ctx context.Context, cur *Schedule, at time.Time, stows []Stow) (*Schedule, error) {
			started <- struct{}{}
			<-ctx.Done()
			return nil, ctx.Err()
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		events, done := make(chan Event), make(chan error, 1)
		go func() { done <- r.Run(ctx, events) }()
		for k := range 3 {
			select {
			case events <- Event{At: takeOverDay.Add(time.Duration(k) * time.Second)}:
			case <-g.wait():
				t.Fatalf("Run did not take event %d", k+1)
			}
			select {
			case <-started:
			case <-g.wait():
				t.Fatalf("the run for event %d never started", k+1)
			}
		}
		cancel()
		select {
		case err := <-done:
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Run(a canceled context) = %v, want %v", err, context.Canceled)
			}
		case <-g.wait():
			t.Fatalf("Run did not return after its context ended")
		}
		if left := goroutinesLeft(g, before); left != 0 {
			t.Errorf("goroutines left after cancel = %d, want 0", left)
		}
	})
}

package main

import (
	"context"
	"errors"
	"runtime"
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

// takeOverDay is midnight UTC on Thursday 28 June 2029.
var takeOverDay = time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)

// at returns h:m on take-over day; hours past 24 run into later days.
func at(h, m int) time.Time {
	return takeOverDay.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
}

// station returns n's site with the code.
func station(n *Network, code string) Site {
	for _, s := range n.Sites {
		if s.Code == code {
			return s
		}
	}
	panic("no station " + code)
}

// useModel puts m in passModel's place until the test ends.
func useModel(t *testing.T, m Predictor) {
	before := passModel
	passModel = m
	t.Cleanup(func() { passModel = before })
}

// checkContacts fails the test unless got holds want's contacts, compared by pass ID and dish.
func checkContacts(t *testing.T, call string, got, want []Contact) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s = %d contacts, want %d", call, len(got), len(want))
	}
	for i := range min(len(got), len(want)) {
		if got[i].ID != want[i].ID || got[i].Dish != want[i].Dish {
			t.Errorf("%s[%d] = %s dish %d, want %s dish %d", call, i, got[i].ID, got[i].Dish, want[i].ID, want[i].Dish)
			return
		}
	}
}

// planFor runs plan in a goroutine and waits up to five seconds for it: real predictions, far more time than they take.
func planFor(t *testing.T, call string, plan func() (*Schedule, error)) (*Schedule, error) {
	t.Helper()
	type result struct {
		s   *Schedule
		err error
	}
	done := make(chan result, 1)
	go func() {
		s, err := plan()
		done <- result{s, err}
	}()
	select {
	case r := <-done:
		return r.s, r.err
	case <-time.After(5 * time.Second):
		t.Fatalf("%s did not return", call)
		return nil, nil
	}
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

func TestPlanUnderLoad(t *testing.T) {
	n := newNetwork()
	cfg := Config{Slots: 12, Days: 7, MinElevDeg: 10, Turnaround: 2 * time.Minute}
	free := handoverRef(oldPlanner(takeOverDay.AddDate(0, 0, -1)), takeOverDay)
	week := Window{From: takeOverDay, To: at(7*24, 0)}
	t.Run("the week from the hand-over, 28 June to 5 July", func(t *testing.T) {
		var m meter // the stations asking passModel at once
		useModel(t, func(e Elements, s Site, from, to time.Time) []PassView {
			m.start()
			defer m.stop()
			return predict(e, s, from, to)
		})
		got, err := planFor(t, "Plan(28 June to 5 July)", func() (*Schedule, error) {
			return Plan(context.Background(), n, week, free, nil, cfg)
		})
		if err != nil {
			t.Fatalf("Plan(28 June to 5 July) error = %v, want <nil>", err)
		}
		if len(got.Contacts) != 14282 || got.Dropped != 28571 {
			t.Errorf("Plan(28 June to 5 July) = %d contacts and %d dropped, want 14282 and 28571", len(got.Contacts), got.Dropped)
		}
		if most := m.most.Load(); most > 12 {
			t.Errorf("most stations asking passModel at once = %d, want 12 at most", most)
		}
	})
	g := newGuard() // from here, every wait is for what a right Plan does at once
	t.Run("a run canceled when Monday 2 July is first asked for", func(t *testing.T) {
		before := runtime.NumGoroutine()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		monday := at(4*24, 0)
		var mu sync.Mutex
		var last time.Time // the latest day passModel was asked for
		var asker string   // the station that asked for it
		useModel(t, func(e Elements, s Site, from, to time.Time) []PassView {
			mu.Lock()
			if from.After(last) {
				last, asker = from, s.Code
			}
			mu.Unlock()
			if !from.Before(monday) {
				cancel()
			}
			return predict(e, s, from, to)
		})
		_, err := planFor(t, "Plan(a canceled run)", func() (*Schedule, error) {
			return Plan(ctx, n, week, free, nil, cfg)
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Plan(a run canceled on 2 July) error = %v, want %v", err, context.Canceled)
		}
		mu.Lock()
		latest, who := last, asker
		mu.Unlock()
		if latest.After(monday.AddDate(0, 0, 1)) {
			t.Errorf("after the cancel, %s asked passModel for %s, want no day after 3 July: check ctx before each day's prediction", who, latest.Format("2 January"))
		}
		if left := goroutinesLeft(g, before); left != 0 {
			t.Errorf("goroutines left after cancel = %d, want 0", left)
		}
	})
}

package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"sync"
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

// holder is a pass model that holds each station's first question until the test lets every station go.
type holder struct {
	mu      sync.Mutex
	asked   []string      // the stations that have asked, in the order they first did
	arrived chan string   // each station's code, at its first question
	release chan struct{} // closed to let every station go
}

// newHolder returns a holder that holds every station.
func newHolder() *holder {
	return &holder{arrived: make(chan string, 20), release: make(chan struct{})}
}

// model answers as the pass model does, once the station's first question has been let go.
func (h *holder) model(e Elements, s Site, from, to time.Time) []PassView {
	h.mu.Lock()
	first := !slices.Contains(h.asked, s.Code)
	if first {
		h.asked = append(h.asked, s.Code)
	}
	h.mu.Unlock()
	if first {
		h.arrived <- s.Code
		<-h.release
	}
	return predict(e, s, from, to)
}

// holding waits for want stations to be held, and returns how many were.
func (h *holder) holding(g guard, want int) int {
	for held := 0; held < want; held++ {
		select {
		case <-h.arrived:
		case <-g.wait():
			return held
		}
	}
	return want
}

func TestPlan(t *testing.T) {
	n := newNetwork()
	cfg := Config{Slots: 12, Days: 1, MinElevDeg: 10, Turnaround: 2 * time.Minute}
	day := Window{From: takeOverDay, To: at(24, 0)}
	var want *Schedule // the answer key, planned once the first check has its twelve stations
	if !t.Run("twelve stations planning at once, never thirteen", func(t *testing.T) {
		g := newGuard() // every wait here is for what a right Plan does at once
		h := newHolder()
		useModel(t, h.model)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var got *Schedule
		var err error
		ran := make(chan struct{})
		go func() {
			defer close(ran)
			got, err = Plan(ctx, n, day, nil, nil, cfg)
		}()
		finish := func() {
			close(h.release)
			select {
			case <-ran:
			case <-time.After(5 * time.Second): // real predictions: far more time than they take
				t.Fatalf("Plan(28 June) did not return once every station was let go")
			}
		}
		if held := h.holding(g, 12); held < 12 {
			cancel() // the run is failed already: let it stop
			finish()
			t.Fatalf("stations planning at once = %d, want 12: plan cfg.Slots stations at a time", held)
		}
		// With twelve stations held, a right Plan starts no thirteenth: 50 ms gives a wrong one time to.
		select {
		case code := <-h.arrived:
			t.Errorf("stations planning at once = 13, want 12 at most: %s started while twelve were held", code)
		case <-time.After(50 * time.Millisecond):
		}
		finish()
		if err != nil {
			t.Fatalf("Plan(28 June) error = %v, want <nil>", err)
		}
		want = planRef(n, day, nil, nil)
		checkContacts(t, "Plan(28 June)", got.Contacts, want.Contacts)
	}) {
		return // the checks below come once Plan runs twelve stations at once
	}
	t.Run("a run canceled while twelve stations plan", func(t *testing.T) {
		g := newGuard() // every wait here is for what a right Plan does at once
		h := newHolder()
		useModel(t, h.model)
		before := runtime.NumGoroutine()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var err error
		ran := make(chan struct{})
		go func() {
			defer close(ran)
			_, err = Plan(ctx, n, day, nil, nil, cfg)
		}()
		held := h.holding(g, 12)
		cancel()
		close(h.release)
		select {
		case <-ran:
		case <-time.After(5 * time.Second): // real predictions: far more time than they take
			t.Fatalf("Plan(a canceled run) did not return")
		}
		if held < 12 {
			t.Fatalf("stations planning at once = %d, want 12: plan cfg.Slots stations at a time", held)
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Plan(a canceled run) error = %v, want %v", err, context.Canceled)
		}
		h.mu.Lock()
		asked := slices.Clone(h.asked)
		h.mu.Unlock()
		if len(asked) != 12 {
			t.Errorf("stations that asked passModel = %d, want 12: once the run is canceled, no other station starts", len(asked))
		}
		if left := goroutinesLeft(g, before); left != 0 {
			t.Errorf("goroutines left after cancel = %d, want 0", left)
		}
	})
	for _, slots := range []int{12, 1} {
		t.Run(fmt.Sprintf("take-over day from scratch, %d slots", slots), func(t *testing.T) {
			c := cfg
			c.Slots = slots
			call := fmt.Sprintf("Plan(28 June, %d slots)", slots)
			got, err := planFor(t, call, func() (*Schedule, error) { return Plan(context.Background(), n, day, nil, nil, c) })
			if err != nil {
				t.Fatalf("%s error = %v, want <nil>", call, err)
			}
			if got.Window != day || got.Dropped != want.Dropped {
				t.Errorf("%s = window from %s, %d dropped, want from 28 June 00:00, %d dropped", call, got.From.Format("2 January 15:04"), got.Dropped, want.Dropped)
			}
			checkContacts(t, call, got.Contacts, want.Contacts)
		})
	}
}

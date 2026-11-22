package main

import (
	"fmt"
	"iter"
	"runtime"
	"slices"
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

// watch records how far a test's source got: whether its loop started, how many values it
// yielded, and whether it ended.
type watch struct {
	started, ended bool
	yielded        int
}

// source yields vals one at a time and records in w that its loop started, how many values
// it yielded and, once it returns, that it ended, as a finished or stopped sequence does.
func source[V any](vals []V, w *watch) iter.Seq[V] {
	return func(yield func(V) bool) {
		w.started = true
		defer func() { w.ended = true }()
		for _, v := range vals {
			w.yielded++
			if !yield(v) {
				return
			}
		}
	}
}

// upTo ranges over seq until it ends or has given limit values, and returns them with
// whatever the loop panicked with.
func upTo[V any](seq iter.Seq[V], limit int) (got []V, panicked any) {
	defer func() { panicked = recover() }()
	for v := range seq {
		got = append(got, v)
		if len(got) == limit {
			break
		}
	}
	return got, nil
}

// report ranges over Compare(old, new) in a goroutine of its own, as upTo does, and returns
// what upTo returned. If the loop is still running when the guard ends, or once it has
// taken 64 MB of memory, a thousand times what a day's walk needs, stuck says so instead.
func report(g guard, old, new iter.Seq[Contact], limit int) (got []Change, panicked any, stuck string) {
	type result struct {
		got      []Change
		panicked any
	}
	results := make(chan result, 1)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	start := mem.TotalAlloc
	go func() {
		got, panicked := upTo(Compare(old, new), limit)
		results <- result{got, panicked}
	}()
	timeout, poll := g.wait(), time.NewTicker(time.Millisecond)
	defer poll.Stop()
	for {
		select {
		case r := <-results:
			return r.got, r.panicked, ""
		case <-timeout:
			return nil, nil, "after a second"
		case <-poll.C:
			if runtime.ReadMemStats(&mem); mem.TotalAlloc-start > 64<<20 {
				return nil, nil, "after taking 64 MB"
			}
		}
	}
}

// show is a change as the messages print it.
func show(c Change) string {
	return fmt.Sprintf("%s %s dish %d", c.Kind, c.Contact.ID, c.Contact.Dish)
}

// checkChanges fails the test unless got holds exactly the changes want shows, in order.
func checkChanges(t *testing.T, call string, got []Change, want []string) {
	t.Helper()
	for i := range max(len(got), len(want)) {
		g, w := "nothing", "nothing"
		if i < len(got) {
			g = show(got[i])
		}
		if i < len(want) {
			w = want[i]
		}
		if g != w {
			t.Fatalf("%s change %d = %s, want %s (%d changes, want %d)", call, i, g, w, len(got), len(want))
		}
	}
}

// index returns the position in cs of the contact with this pass ID.
func index(cs []Contact, id string) int {
	return slices.IndexFunc(cs, func(c Contact) bool { return c.ID == id })
}

// without returns a copy of cs without the contact with this pass ID.
func without(cs []Contact, id string) []Contact {
	return slices.Delete(slices.Clone(cs), index(cs, id), index(cs, id)+1)
}

// swapped returns a copy of cs with the contact with this pass ID and the one after it
// in each other's places.
func swapped(cs []Contact, id string) []Contact {
	out, i := slices.Clone(cs), index(cs, id)
	out[i], out[i+1] = out[i+1], out[i]
	return out
}

// onDish returns a copy of cs with the contact with this pass ID on another dish.
func onDish(cs []Contact, id string, dish int) []Contact {
	out := slices.Clone(cs)
	out[index(cs, id)].Dish = dish
	return out
}

func TestCompare(t *testing.T) {
	oldPlan, newPlan := thursday()
	// Both rise at 05:00:11; the old planner lists Brine-3's first.
	sqt, vsk := "brine-3/SQT/16651", "furrow-a/VSK/23251"
	last := oldPlan[len(oldPlan)-1].ID // brine-3/BRM/16663, the day's last contact

	tests := []struct {
		name     string
		old, new []Contact
		want     []string
	}{
		{"Thursday against itself, listed both ways at 05:00:11", oldPlan, swapped(oldPlan, sqt), nil},
		{"Thursday from both planners", oldPlan, newPlan, []string{
			"new only furrow-b/TIR/23250 dish 1",
			"new only pipit-2/TIR/21175 dish 1",
			"new only pipit-3/TIR/11124 dish 1",
			"dish furrow-a/TIR/23253 dish 2",
		}},
		{"Brine-3 over Salar Quieto missing from passplan 2's", oldPlan, without(oldPlan, sqt), []string{"old only brine-3/SQT/16651 dish 1"}},
		{"Brine-3 missing from passplan 2's, 05:00:11 in the other order", swapped(oldPlan, sqt), without(oldPlan, sqt), []string{"old only brine-3/SQT/16651 dish 1"}},
		{"Furrow-A over Varsk Island missing from the old planner's", without(oldPlan, vsk), swapped(oldPlan, sqt), []string{"new only furrow-a/VSK/23251 dish 1"}},
		{"Furrow-A missing from the old planner's, 05:00:11 in the other order", without(oldPlan, vsk), oldPlan, []string{"new only furrow-a/VSK/23251 dish 1"}},
		{"Brine-3 over Salar Quieto on dish 2 in passplan 2's", oldPlan, onDish(swapped(oldPlan, sqt), sqt, 2), []string{"dish brine-3/SQT/16651 dish 2"}},
		{"the day's last contact missing from passplan 2's", oldPlan, without(oldPlan, last), []string{"old only brine-3/BRM/16663 dish 1"}},
		{"the day's last contact missing from the old planner's", without(oldPlan, last), oldPlan, []string{"new only brine-3/BRM/16663 dish 1"}},
	}
	g := newGuard() // for Compare alone: the schedules above are planned already
	stuck := false
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, panicked, still := report(g, slices.Values(tc.old), slices.Values(tc.new), len(tc.want)+10)
			if still != "" {
				stuck = true
				t.Fatalf("Compare(old, new) was still running %s, with %d contacts against %d", still, len(tc.old), len(tc.new))
			}
			if panicked != nil {
				t.Fatalf("ranging over Compare(old, new) panicked: %v", panicked)
			}
			checkChanges(t, "Compare(old, new)", got, tc.want)
		})
		if stuck {
			return // that walk runs on beside any later case, so stop here
		}
	}

	t.Run("a report that stops at the first change", func(t *testing.T) {
		var a, b watch
		got, panicked, still := report(g, source(oldPlan, &a), source(newPlan, &b), 1)
		if still != "" {
			t.Fatalf("Compare(Thursday's schedules) was still running %s", still)
		}
		if panicked != nil {
			t.Fatalf("breaking after the first change panicked: %v", panicked)
		}
		checkChanges(t, "Compare(Thursday's schedules), for a loop that stops at 1,", got, []string{"new only furrow-b/TIR/23250 dish 1"})
		for _, s := range []struct {
			name string
			w    watch
			all  int
		}{{"old", a, len(oldPlan)}, {"new", b, len(newPlan)}} {
			if s.w.started && !s.w.ended {
				t.Errorf("after a loop over Compare(old, new) stopped at its first change, %s's loop had not ended, want both stopped", s.name)
			}
			if s.w.yielded == s.all {
				t.Errorf("Compare(old, new) read all %d of %s's contacts for a loop that stops at its first change, at 07:38:06, want it to pull no further than that", s.all, s.name)
			}
		}
	})
}

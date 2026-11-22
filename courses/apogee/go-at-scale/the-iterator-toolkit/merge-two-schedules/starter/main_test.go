package main

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"testing"
)

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

// contactsAt returns the contacts of the station with this code, in the schedule's order.
func contactsAt(sch *Schedule, code string) []Contact {
	var out []Contact
	for _, c := range sch.Contacts {
		if c.Station == code {
			out = append(out, c)
		}
	}
	return out
}

// contactsByAOS orders contacts as byAOS orders their passes.
func contactsByAOS(a, b Contact) int { return byAOS(a.PassView, b.PassView) }

// tagged is a contact with the planner that planned it.
type tagged struct {
	planner string
	id      string
	c       Contact
}

// tag returns cs tagged with planner.
func tag(cs []Contact, planner string) []tagged {
	var out []tagged
	for _, c := range cs {
		out = append(out, tagged{planner, c.ID, c})
	}
	return out
}

// checkStopped fails the test if a source's loop started and never ended.
func checkStopped(t *testing.T, call, name string, w watch) {
	t.Helper()
	if w.started && !w.ended {
		t.Errorf("after %s, %s's loop had not ended, want it stopped on every way out", call, name)
	}
}

func TestMerge(t *testing.T) {
	monday := firstWaveDay()
	tir, vsk := contactsAt(monday, "TIR"), contactsAt(monday, "VSK")

	t.Run("Cape Wren's gusts and Ostrel Island's", func(t *testing.T) {
		cwr, ost := []int{58, 64, 72, 81}, []int{61, 64, 70, 77, 90}
		want := []int{58, 61, 64, 64, 70, 72, 77, 81, 90}
		got, panicked := upTo(Merge(source(cwr, new(watch)), source(ost, new(watch)), cmp.Compare[int]), len(want)+1)
		if panicked != nil {
			t.Fatalf("ranging over Merge(%v, %v) panicked: %v", cwr, ost, panicked)
		}
		if !slices.Equal(got, want) {
			t.Errorf("Merge(%v, %v) = %v, want %v", cwr, ost, got, want)
		}
	})
	t.Run("Tirrin Rock and Varsk Island on Monday", func(t *testing.T) {
		call := "Merge(Tirrin Rock's contacts, Varsk Island's, by AOS)"
		got, panicked := upTo(Merge(source(tir, new(watch)), source(vsk, new(watch)), contactsByAOS), len(tir)+len(vsk)+1)
		if panicked != nil {
			t.Fatalf("ranging over %s panicked: %v", call, panicked)
		}
		if len(got) != len(tir)+len(vsk) {
			t.Fatalf("%s yielded %d contacts, want %d", call, len(got), len(tir)+len(vsk))
		}
		var gotIDs, wantIDs []string
		for i, c := range got {
			if i > 0 && c.AOS.Before(got[i-1].AOS) {
				t.Fatalf("%s yielded %s before %s, want AOS order", call, got[i-1].ID, c.ID)
			}
			gotIDs = append(gotIDs, c.ID)
		}
		for _, c := range append(slices.Clone(tir), vsk...) {
			wantIDs = append(wantIDs, c.ID)
		}
		slices.Sort(gotIDs)
		slices.Sort(wantIDs)
		if !slices.Equal(gotIDs, wantIDs) {
			t.Errorf("%s did not yield every contact of both stations once", call)
		}
	})
	t.Run("both planners' Tirrin Rock, the old planner's first", func(t *testing.T) {
		call := "Merge(the old planner's TIR contacts, passplan 2's, by AOS)"
		taggedByAOS := func(a, b tagged) int { return contactsByAOS(a.c, b.c) }
		got, panicked := upTo(Merge(source(tag(tir, "the old planner"), new(watch)), source(tag(tir, "passplan 2"), new(watch)), taggedByAOS), 2*len(tir)+1)
		if panicked != nil {
			t.Fatalf("ranging over %s panicked: %v", call, panicked)
		}
		if len(got) != 2*len(tir) {
			t.Fatalf("%s yielded %d contacts, want %d", call, len(got), 2*len(tir))
		}
		for i, g := range got {
			want := tagged{"the old planner", tir[i/2].ID, tir[i/2]}
			if i%2 == 1 {
				want.planner = "passplan 2"
			}
			if g.planner != want.planner || g.id != want.id {
				t.Fatalf("%s value %d = %s's %s, want %s's %s", call, i, g.planner, g.id, want.planner, want.id)
			}
		}
	})
	t.Run("an empty schedule", func(t *testing.T) {
		for _, order := range []string{"first", "second"} {
			a, b := tir, []Contact{}
			if order == "second" {
				a, b = b, a
			}
			call := "Merge with an empty schedule " + order
			got, panicked := upTo(Merge(source(a, new(watch)), source(b, new(watch)), contactsByAOS), len(tir)+1)
			if panicked != nil {
				t.Fatalf("ranging over %s panicked: %v", call, panicked)
			}
			if len(got) != len(tir) {
				t.Fatalf("%s yielded %d contacts, want Tirrin Rock's %d", call, len(got), len(tir))
			}
			for i := range got {
				if got[i].ID != tir[i].ID {
					t.Fatalf("%s contact %d = %s, want %s", call, i, got[i].ID, tir[i].ID)
				}
			}
		}
	})
	t.Run("logs that stop after 1 to 12 contacts", func(t *testing.T) {
		for limit := 1; limit <= 12; limit++ {
			var a, b watch
			call := fmt.Sprintf("a loop over Merge(TIR, VSK) that stops at %d", limit)
			got, panicked := upTo(Merge(source(tir, &a), source(vsk, &b), contactsByAOS), limit)
			if panicked != nil {
				t.Fatalf("breaking after %d contacts panicked: %v", limit, panicked)
			}
			if len(got) != limit {
				t.Fatalf("Merge(TIR, VSK) gave %d contacts to a loop that stops at %d", len(got), limit)
			}
			checkStopped(t, call, "a", a)
			checkStopped(t, call, "b", b)
			if a.yielded == len(tir) || b.yielded == len(vsk) {
				t.Errorf("%s read %d of a's %d contacts and %d of b's %d, want it to pull only as it yields", call, a.yielded, len(tir), b.yielded, len(vsk))
			}
			if t.Failed() {
				return // one stopping point is enough to show the mistake
			}
		}
	})
}

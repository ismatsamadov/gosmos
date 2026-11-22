package main

import (
	"iter"
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

// row is one row of the test Board: the pass IDs of two contacts side by side.
type row struct{ left, right string }

// rows ranges over seq until it ends or has given limit pairs, and returns them as rows
// with whatever the loop panicked with.
func rows(seq iter.Seq2[Contact, Contact], limit int) (got []row, panicked any) {
	defer func() { panicked = recover() }()
	for a, b := range seq {
		got = append(got, row{a.ID, b.ID})
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

// checkRows fails the test unless got pairs left[i] with right[i] for each of the first n.
func checkRows(t *testing.T, call string, got []row, left, right []Contact, n int) {
	t.Helper()
	if len(got) != n {
		t.Fatalf("%s yielded %d pairs, want %d", call, len(got), n)
	}
	for i, r := range got {
		if want := (row{left[i].ID, right[i].ID}); r != want {
			t.Fatalf("%s pair %d = (%s, %s), want (%s, %s)", call, i, r.left, r.right, want.left, want.right)
		}
	}
}

// checkStopped fails the test if a source's loop started and never ended.
func checkStopped(t *testing.T, call, name string, w watch) {
	t.Helper()
	if w.started && !w.ended {
		t.Errorf("after %s, %s's loop had not ended, want it stopped on every way out", call, name)
	}
}

func TestZip(t *testing.T) {
	oldPlan, newPlan := firstWaveDay(), firstWaveDay() // on Monday the two planners agree
	oldTIR, newTIR := contactsAt(oldPlan, "TIR"), contactsAt(newPlan, "TIR")

	t.Run("five rows beside seven", func(t *testing.T) {
		var a, b watch
		call := "Zip(5 of the old planner's contacts, 7 of passplan 2's)"
		got, panicked := rows(Zip(source(oldTIR[:5], &a), source(newTIR[:7], &b)), 8)
		if panicked != nil {
			t.Fatalf("ranging over %s panicked: %v", call, panicked)
		}
		checkRows(t, call, got, oldTIR, newTIR, 5)
		checkStopped(t, call, "a", a)
		checkStopped(t, call, "b", b)
	})
	t.Run("seven rows beside five", func(t *testing.T) {
		var a, b watch
		call := "Zip(7 of the old planner's contacts, 5 of passplan 2's)"
		got, panicked := rows(Zip(source(oldTIR[:7], &a), source(newTIR[:5], &b)), 8)
		if panicked != nil {
			t.Fatalf("ranging over %s panicked: %v", call, panicked)
		}
		checkRows(t, call, got, oldTIR, newTIR, 5)
		checkStopped(t, call, "a", a)
		checkStopped(t, call, "b", b)
	})
	t.Run("Tirrin Rock beside its partner Ystmark", func(t *testing.T) {
		yst := contactsAt(newPlan, "YST")
		call := "Zip(Tirrin Rock's contacts, Ystmark's)"
		got, panicked := rows(Zip(source(oldTIR, new(watch)), source(yst, new(watch))), len(oldTIR)+1)
		if panicked != nil {
			t.Fatalf("ranging over %s panicked: %v", call, panicked)
		}
		checkRows(t, call, got, oldTIR, yst, min(len(oldTIR), len(yst)))
	})
	t.Run("an empty schedule beside seven", func(t *testing.T) {
		var a, b watch
		call := "Zip(no contacts, 7 contacts)"
		got, panicked := rows(Zip(source([]Contact{}, &a), source(newTIR[:7], &b)), 1)
		if panicked != nil {
			t.Fatalf("ranging over %s panicked: %v", call, panicked)
		}
		checkRows(t, call, got, nil, nil, 0)
		checkStopped(t, call, "b", b)
	})
	t.Run("seven beside an empty schedule", func(t *testing.T) {
		var a, b watch
		call := "Zip(7 contacts, no contacts)"
		got, panicked := rows(Zip(source(oldTIR[:7], &a), source([]Contact{}, &b)), 1)
		if panicked != nil {
			t.Fatalf("ranging over %s panicked: %v", call, panicked)
		}
		checkRows(t, call, got, nil, nil, 0)
		checkStopped(t, call, "a", a)
		checkStopped(t, call, "b", b)
	})
	t.Run("a test Board with room for two rows", func(t *testing.T) {
		var a, b watch
		call := "Zip(both planners' TIR contacts) for a loop that stops at 2"
		got, panicked := rows(Zip(source(oldTIR, &a), source(newTIR, &b)), 2)
		if panicked != nil {
			t.Fatalf("breaking after 2 pairs panicked: %v", panicked)
		}
		checkRows(t, call, got, oldTIR, newTIR, 2)
		checkStopped(t, call, "a", a)
		checkStopped(t, call, "b", b)
		if b.yielded == len(newTIR) {
			t.Errorf("%s read all %d of b's contacts, want only the 2 its rows need", call, len(newTIR))
		}
	})
}

package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
	"time"
)

// pass returns a predicted pass from its ID, its AOS and LOS ("2029-04-08 10:29:34", UTC) and its max elevation.
func pass(id, aos, los string, elevDeg float64) PassView {
	sat, rest, _ := strings.Cut(id, "/")
	station, _, _ := strings.Cut(rest, "/")
	at := func(s string) time.Time {
		t, _ := time.Parse(time.DateTime, s)
		return t
	}
	return PassView{ID: id, Pass: Pass{SatID: sat, Station: station, AOS: at(aos), LOS: at(los), MaxElevDeg: elevDeg}}
}

// portEllery is Port Ellery's passes on Sunday 8 April from 09:53 to 11:35, by AOS and then ID.
var portEllery = []PassView{
	pass("furrow-b/PEL/23297", "2029-04-08 09:53:24", "2029-04-08 10:01:12", 19.4),
	pass("brine-3/PEL/16701", "2029-04-08 10:29:34", "2029-04-08 10:38:49", 67.7),
	pass("pipit-3/PEL/11171", "2029-04-08 10:29:34", "2029-04-08 10:35:05", 10.1),
	pass("furrow-a/PEL/23300", "2029-04-08 10:39:14", "2029-04-08 10:48:28", 86.6),
	pass("cirrine-1/PEL/18660", "2029-04-08 10:45:40", "2029-04-08 10:54:50", 25.7),
	pass("furrow-b/PEL/23298", "2029-04-08 11:26:48", "2029-04-08 11:35:28", 33.1),
}

// safePop pops q, turning a panic into an error, so a test can report it.
func safePop[T any](q *Queue[T]) (v T, ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	v, ok = q.Pop()
	return v, ok, nil
}

func TestQueue(t *testing.T) {
	t.Run("Sunday's passes at Port Ellery, by byAOS", func(t *testing.T) {
		q := NewQueue(byAOS)
		for _, i := range []int{3, 5, 2, 0, 4, 1} {
			q.Push(portEllery[i])
		}
		if q.Len() != len(portEllery) {
			t.Fatalf("Len() after %d pushes = %d, want %d", len(portEllery), q.Len(), len(portEllery))
		}
		for i, want := range portEllery {
			if p, ok := q.Pop(); !ok || p.ID != want.ID {
				t.Fatalf("Pop() number %d = %q, %v, want %q, true", i+1, p.ID, ok, want.ID)
			}
		}
		if q.Len() != 0 {
			t.Errorf("Len() after popping every pass = %d, want 0", q.Len())
		}
	})
	t.Run("Cape Wren's gusts, ints by cmp.Compare", func(t *testing.T) {
		q := NewQueue(cmp.Compare[int])
		for _, kmh := range []int{72, 50, 81, 64} {
			q.Push(kmh)
		}
		for i, want := range []int{50, 64, 72, 81} {
			if got, ok := q.Pop(); !ok || got != want {
				t.Fatalf("Pop() number %d = %d, %v, want %d, true", i+1, got, ok, want)
			}
		}
	})
	t.Run("an empty queue", func(t *testing.T) {
		passes := NewQueue(byAOS)
		if p, ok, err := safePop(passes); err != nil {
			t.Errorf("Pop() on an empty queue of passes panicked: %v", err)
		} else if ok || p.ID != "" {
			t.Errorf("Pop() on an empty queue of passes = %q, %v, want the zero PassView and false", p.ID, ok)
		}
		gusts := NewQueue(cmp.Compare[int])
		gusts.Push(72)
		gusts.Pop()
		if got, ok, err := safePop(gusts); err != nil {
			t.Errorf("Pop() on a queue its last Pop emptied panicked: %v", err)
		} else if ok || got != 0 {
			t.Errorf("Pop() on a queue its last Pop emptied = %d, %v, want 0, false", got, ok)
		}
	})
	t.Run("1,000 ints from a seeded source", func(t *testing.T) {
		r := rand.New(rand.NewPCG(1, 2))
		q := NewQueue(cmp.Compare[int])
		want := make([]int, 1000)
		for i := range want {
			want[i] = r.IntN(100_000)
			q.Push(want[i])
		}
		slices.Sort(want)
		for i, w := range want {
			if got, ok := q.Pop(); !ok || got != w {
				t.Fatalf("Pop() number %d = %d, %v, want %d, true", i+1, got, ok, w)
			}
		}
	})
}

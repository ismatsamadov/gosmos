package main

import (
	"iter"
	"slices"
	"testing"
)

// counted yields vals one at a time, as a schedule does, and counts in *n the values it has
// yielded.
func counted[V any](vals []V, n *int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			*n++
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

// ids returns the pass IDs of cs.
func ids(cs []Contact) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.ID)
	}
	return out
}

// at keeps the contacts of the station with this code.
func at(code string) func(Contact) bool {
	return func(c Contact) bool { return c.Station == code }
}

func TestFilter(t *testing.T) {
	monday := firstWaveDay()
	var tir []Contact
	for _, c := range monday.Contacts {
		if c.Station == "TIR" {
			tir = append(tir, c)
		}
	}

	t.Run("Tirrin Rock's panel on the Board", func(t *testing.T) {
		got, panicked := upTo(Filter(slices.Values(monday.Contacts), at("TIR")), len(tir)+1)
		if panicked != nil {
			t.Fatalf("ranging over Filter(Monday's contacts, at TIR) panicked: %v", panicked)
		}
		if len(got) > len(tir) {
			t.Fatalf("Filter(Monday's %d contacts, at TIR) yielded more than %d contacts, want Tirrin Rock's %d", len(monday.Contacts), len(tir), len(tir))
		}
		if len(got) != len(tir) {
			t.Fatalf("Filter(Monday's %d contacts, at TIR) yielded %d contacts, want Tirrin Rock's %d", len(monday.Contacts), len(got), len(tir))
		}
		for i := range got {
			if got[i].ID != tir[i].ID {
				t.Fatalf("Filter(Monday's contacts, at TIR) contact %d = %s, want %s", i, got[i].ID, tir[i].ID)
			}
		}
	})
	t.Run("Ilvane Spit, live from 7 May", func(t *testing.T) {
		got, panicked := upTo(Filter(slices.Values(monday.Contacts), at("ILV")), 1)
		if panicked != nil {
			t.Fatalf("ranging over Filter(Monday's contacts, at ILV) panicked: %v", panicked)
		}
		if len(got) != 0 {
			t.Errorf("Filter(Monday's contacts, at ILV) yielded %s, want nothing", got[0].ID)
		}
	})
	t.Run("a panel with room for two", func(t *testing.T) {
		asked := 0
		got, panicked := upTo(Filter(counted(monday.Contacts, &asked), at("TIR")), 2)
		if panicked != nil {
			t.Fatalf("breaking after Tirrin Rock's second contact panicked: %v", panicked)
		}
		if len(got) != 2 || got[0].ID != tir[0].ID || got[1].ID != tir[1].ID {
			t.Fatalf("Filter(Monday's contacts, at TIR) gave %v to a loop that stops at 2, want [%s %s]", ids(got), tir[0].ID, tir[1].ID)
		}
		want := slices.IndexFunc(monday.Contacts, func(c Contact) bool { return c.ID == tir[1].ID }) + 1
		if asked != want {
			t.Errorf("the schedule yielded %d contacts to Filter(at TIR) for a loop that stopped at Tirrin Rock's second, want %d", asked, want)
		}
	})
	t.Run("Tirrin Rock's dish numbers", func(t *testing.T) {
		var dishes []int
		for _, c := range tir {
			dishes = append(dishes, c.Dish)
		}
		got, panicked := upTo(Filter(slices.Values(dishes), func(d int) bool { return d == 2 }), len(dishes)+1)
		if panicked != nil {
			t.Fatalf("ranging over Filter(Tirrin Rock's dishes, dish 2) panicked: %v", panicked)
		}
		twos := 0
		for _, d := range got {
			if d == 2 {
				twos++
			}
		}
		if len(got) != 13 || twos != 13 {
			t.Errorf("Filter(Tirrin Rock's %d dish numbers, dish 2) yielded %d numbers, %d of them 2, want 13, all 2", len(dishes), len(got), twos)
		}
	})
}

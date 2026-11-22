package main

import (
	"iter"
	"testing"
)

// shown is one line of a station panel: a dish and the contact it flies.
type shown struct {
	dish int
	id   string
}

// panel ranges over seq until it ends or has given limit pairs, and returns them with
// whatever the loop panicked with.
func panel(seq iter.Seq2[int, Contact], limit int) (got []shown, panicked any) {
	defer func() { panicked = recover() }()
	for dish, c := range seq {
		got = append(got, shown{dish, c.ID})
		if len(got) == limit {
			break
		}
	}
	return got, nil
}

func TestAt(t *testing.T) {
	sch := firstWaveDay()
	var want []shown
	for _, c := range sch.Contacts {
		if c.Station == "TIR" {
			want = append(want, shown{c.Dish, c.ID})
		}
	}

	t.Run("Tirrin Rock's first day", func(t *testing.T) {
		got, panicked := panel(sch.At("TIR"), len(want)+1)
		if panicked != nil {
			t.Fatalf("ranging over At(TIR) panicked: %v", panicked)
		}
		if len(got) > len(want) {
			t.Fatalf("At(TIR) yielded more than %d pairs, want %d: Tirrin Rock's contacts and no others", len(want), len(want))
		}
		if len(got) != len(want) {
			t.Fatalf("At(TIR) yielded %d pairs, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("At(TIR) pair %d = dish %d, %s, want dish %d, %s", i, got[i].dish, got[i].id, want[i].dish, want[i].id)
			}
		}
	})
	t.Run("Ilvane Spit, live from 7 May", func(t *testing.T) {
		if got, _ := panel(sch.At("ILV"), 1); len(got) != 0 {
			t.Errorf("At(ILV) yielded dish %d, %s, want nothing", got[0].dish, got[0].id)
		}
	})
	t.Run("a panel with room for two", func(t *testing.T) {
		got, panicked := panel(sch.At("TIR"), 2)
		if panicked != nil {
			t.Fatalf("breaking after 2 pairs panicked: %v", panicked)
		}
		if len(got) != 2 {
			t.Fatalf("At(TIR) gave %d pairs to a loop that stops at 2, want 2", len(got))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("At(TIR) pair %d = dish %d, %s, want dish %d, %s", i, got[i].dish, got[i].id, want[i].dish, want[i].id)
			}
		}
	})
	t.Run("two panels, one schedule", func(t *testing.T) {
		seq := sch.At("TIR")
		first, _ := panel(seq, len(want)+1)
		second, _ := panel(seq, len(want)+1)
		if len(first) != len(want) || len(second) != len(want) {
			t.Fatalf("ranging twice over one At(TIR) gave %d pairs, then %d, want %d both times", len(first), len(second), len(want))
		}
		for i := range want {
			if first[i] != want[i] || second[i] != want[i] {
				t.Fatalf("ranging twice over one At(TIR): pair %d = dish %d, %s, then dish %d, %s, want dish %d, %s both times", i, first[i].dish, first[i].id, second[i].dish, second[i].id, want[i].dish, want[i].id)
			}
		}
	})
}

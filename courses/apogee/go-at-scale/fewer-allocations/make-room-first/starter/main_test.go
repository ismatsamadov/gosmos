package main

import (
	"slices"
	"testing"
	"time"
)

func TestCollect(t *testing.T) {
	_, cwr := stationDay("CWR", june28)
	may21 := time.Date(2029, 5, 21, 0, 0, 0, 0, time.UTC)
	_, vrl := stationDay("VRL", may21)
	tests := []struct {
		name   string
		passes []PassView
		n      int
	}{
		{"Cape Wren's day, as many as last run", cwr, len(cwr)},
		{"Cape Wren's day, room to spare", cwr, 512},
		{"Cape Wren's day, more than last run", cwr, 300},
		{"Cape Wren's day, a first run with no count", cwr, 0},
		{"Vorrel Island before it goes live", vrl, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collect(slices.Values(tt.passes), tt.n)
			if len(got) != len(tt.passes) {
				t.Fatalf("len(collect(seq, %d)) = %d, want %d: every pass, once", tt.n, len(got), len(tt.passes))
			}
			for i, p := range tt.passes {
				if got[i] != p {
					t.Fatalf("collect(seq, %d)[%d] = %q, want %q", tt.n, i, got[i].ID, p.ID)
				}
			}
		})
	}
}

func TestCollectAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	_, cwr := stationDay("CWR", june28)
	var got []PassView
	// What collect costs with no passes and no room to make:
	// nothing when the compiler copies it into this test, a few
	// for its loop when it can't. Only allocations past those
	// count.
	base := testing.AllocsPerRun(10, func() {
		got = collect(slices.Values(cwr[:0]), 0)
	})
	allocs := testing.AllocsPerRun(10, func() {
		got = collect(slices.Values(cwr), len(cwr))
	}) - base
	if len(got) != len(cwr) {
		return // TestCollect says what is wrong
	}
	if allocs != 1 {
		t.Errorf("collect(Cape Wren's %d candidates, %d) made %.0f allocations, want 1: room for all %d, made once", len(cwr), len(cwr), allocs, len(cwr))
	}
}

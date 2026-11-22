package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestAccept(t *testing.T) {
	for _, c := range []struct {
		name            string
		arrive          []uint32
		shown           []uint32
		missed, dropped int
	}{
		{"five in order", []uint32{52817, 52818, 52819, 52820, 52821}, []uint32{52817, 52818, 52819, 52820, 52821}, 0, 0},
		{"one that came twice", []uint32{68841, 68842, 68842, 68843}, []uint32{68841, 68842, 68843}, 0, 1},
		{"two that changed places", []uint32{68841, 68843, 68842, 68844}, []uint32{68841, 68843, 68844}, 1, 1},
		{"three lost on the line", []uint32{68841, 68845}, []uint32{68841, 68845}, 3, 0},
		{"a first frame of any number", []uint32{68841}, []uint32{68841}, 0, 0},
		{"frame 0, twice", []uint32{0, 0, 1}, []uint32{0, 1}, 0, 1},
		{"an old frame, a repeat and a late one", []uint32{68850, 68841, 68850, 68849, 68851}, []uint32{68850, 68851}, 0, 3},
	} {
		t.Run(c.name, func(t *testing.T) {
			var tr Tracker
			var shown []uint32
			for _, seq := range c.arrive {
				if tr.Accept(seq) {
					shown = append(shown, seq)
				}
			}
			if !slices.Equal(shown, c.shown) {
				t.Errorf("Accept of %v showed %v, want %v", c.arrive, shown, c.shown)
			}
			want := Tracker{Last: c.shown[len(c.shown)-1], Shown: len(c.shown), Missed: c.missed, Dropped: c.dropped}
			if tr != want {
				t.Errorf("after Accept of %v: Tracker = %s, want %s", c.arrive, show(tr), show(want))
			}
		})
	}
}

// show writes a Tracker's fields.
func show(t Tracker) string {
	return fmt.Sprintf("Last %d, Shown %d, Missed %d, Dropped %d", t.Last, t.Shown, t.Missed, t.Dropped)
}

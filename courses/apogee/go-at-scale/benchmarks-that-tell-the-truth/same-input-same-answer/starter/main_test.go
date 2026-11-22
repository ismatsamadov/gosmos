package main

import (
	"slices"
	"testing"
	"time"
)

func TestSameContacts(t *testing.T) {
	n := newNetwork()
	thursday := time.Date(2029, 5, 17, 0, 0, 0, 0, time.UTC)
	day := Window{From: thursday, To: thursday.AddDate(0, 0, 1)}
	cfg := Config{MinElevDeg: 10, Turnaround: 2 * time.Minute}
	old, new := planOld(n, day, cfg), planNew(n, day, cfg)
	if len(old) != 1231 || len(new) != 1231 {
		t.Fatalf("planOld and planNew planned %d and %d contacts on 17 May, want 1231 each", len(old), len(new))
	}

	dish213 := slices.Clone(new)
	dish213[213].Dish = 2
	tie := slices.Clone(new)
	tie[369], tie[370] = tie[370], tie[369]
	two := slices.Clone(dish213)
	two[369].Dish = 1
	short := new[:len(new)-1]
	long := append([]Contact{new[0]}, new...)

	cases := []struct {
		name     string
		old, new []Contact
		want     string // "" for nil
	}{
		{"17 May from both planners", old, new, ""},
		{"contact 213 on Ilvane Spit's other dish", old, dish213,
			"contact 213: old cirrine-5/ILV/00445 dish 1, new cirrine-5/ILV/00445 dish 2"},
		{"the 07:08:09 tie broken the other way", old, tie,
			"contact 369: old pipit-3/TIR/11762 dish 2, new pipit-3/YST/11762 dish 2"},
		{"two differences, at 213 and 369", old, two,
			"contact 213: old cirrine-5/ILV/00445 dish 1, new cirrine-5/ILV/00445 dish 2"},
		{"new without the day's last contact", old, short, "old has 1231 contacts, new has 1230"},
		{"new with its first contact twice", old, long, "old has 1231 contacts, new has 1232"},
		{"two empty lists", nil, []Contact{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := sameContacts(c.old, c.new)
			switch {
			case c.want == "" && err != nil:
				t.Errorf("sameContacts = %q, want nil", err)
			case c.want != "" && err == nil:
				t.Errorf("sameContacts = nil, want %q", c.want)
			case c.want != "" && err.Error() != c.want:
				t.Errorf("sameContacts = %q, want %q", err, c.want)
			}
		})
	}
}

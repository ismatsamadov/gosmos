package main

import (
	"iter"
	"slices"
	"time"
)

// byAOS orders passes by AOS.
func byAOS(a, b PassView) int {
	return a.AOS.Compare(b.AOS)
}

// canServe reports whether the station offers one of the satellite's bands.
func canServe(s Site, sat Satellite) bool {
	return slices.ContainsFunc(sat.Bands, func(b Band) bool { return slices.Contains(s.Bands, b) })
}

// assign gives one station's candidate passes dishes. The passes come in AOS
// order, then ID; free[d] is when dish d+1 is next free (zero: free now).
// It returns the contacts in AOS order and how many passes found no dish.
func assign(s Site, passes iter.Seq[PassView], free []time.Time, turnaround time.Duration) ([]Contact, int) {
	next := make([]time.Time, s.Dishes)
	copy(next, free)
	var out []Contact
	dropped := 0
	for p := range passes {
		d := slices.IndexFunc(next, func(t time.Time) bool { return !p.AOS.Before(t) })
		if d < 0 {
			dropped++
			continue
		}
		next[d] = p.LOS.Add(turnaround)
		out = append(out, Contact{PassView: p, Dish: d + 1})
	}
	return out, dropped
}

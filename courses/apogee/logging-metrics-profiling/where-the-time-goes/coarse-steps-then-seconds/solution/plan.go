package main

import "time"

// contact is one pass's AOS and LOS as the elevation alone shows them: the moments it crosses the mask,
// before they are rounded to the second.
type contact struct{ AOS, LOS time.Time }

// Planned is one predicted pass in the day's plan.
type Planned struct {
	PassView        // ID and Pass (Season 3)
	Booking  string // "bk-5803", or "" when nobody booked it
}

// findPasses returns el's passes over maskDeg that start and end from from to to, by the rules of
// scanEverySecond.
func findPasses(el func(time.Time) float64, from, to time.Time, maskDeg float64) []contact {
	var passes []contact
	var aos time.Time
	was := el(from) >= maskDeg
	last := to.Add(-time.Second)
	for prev := from; prev.Before(last); {
		t := prev.Add(30 * time.Second)
		if t.After(last) {
			t = last
		}
		if now := el(t) >= maskDeg; now != was {
			at := crossing(el, prev, t, maskDeg)
			if now {
				aos = at
			} else if !aos.IsZero() && toSecond(at).Sub(toSecond(aos)) >= 2*time.Minute {
				passes = append(passes, contact{aos, at})
			}
			was = now
		}
		prev = t
	}
	return passes
}

package main

import (
	"slices"
	"time"
)

// contact is one pass's AOS and LOS as the elevation alone shows them: the moments it crosses the mask,
// before they are rounded to the second.
type contact struct{ AOS, LOS time.Time }

// bookingKey is what a booking and its pass share: satellite, station and AOS in Unix seconds.
type bookingKey struct {
	SatID, Station string
	AOSUnix        int64
}

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

// attach gives each pass its booking: the first in bookings with the pass's satellite, station and
// AOS, or none.
func attach(passes []PassView, bookings []Booking) []Planned {
	index := make(map[bookingKey]string, len(bookings))
	for _, b := range bookings {
		k := bookingKey{b.SatID, b.Station, b.AOS.Unix()}
		if _, ok := index[k]; !ok {
			index[k] = b.ID
		}
	}
	planned := make([]Planned, 0, len(passes))
	for _, p := range passes {
		planned = append(planned, Planned{PassView: p, Booking: index[bookingKey{p.SatID, p.Station, p.AOS.Unix()}]})
	}
	return planned
}

// buildDay joins the groups' passes into one day, ordered by AOS, then ID.
func buildDay(groups [][]Planned) []Planned {
	var day []Planned
	for _, g := range groups {
		day = append(day, g...)
	}
	slices.SortFunc(day, earlier)
	return day
}

// planDay plans every pass of sats over stations from from to to, each with its booking.
func planDay(sats []Satellite, els map[string]Elements, stations []Station, bookings []Booking, from, to time.Time) []Planned {
	var passes []PassView
	for _, sat := range sats {
		elems := els[sat.ID]
		for _, st := range stations {
			el := elevation(sat, elems, st)
			passes = append(passes, passViews(sat, elems, st, findPasses(el, from, to, st.MaskDeg))...)
		}
	}
	return buildDay([][]Planned{attach(passes, bookings)})
}

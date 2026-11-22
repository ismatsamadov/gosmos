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

// planDay plans every pass of sats over stations from from to to, each with its booking.
func planDay(sats []Satellite, els map[string]Elements, stations []Station, bookings []Booking, from, to time.Time) []Planned {
	var groups [][]Planned
	for _, sat := range sats {
		elems := els[sat.ID]
		for _, st := range stations {
			el := elevation(sat, elems, st)
			passes := passViews(sat, elems, st, scanEverySecond(el, from, to, st.MaskDeg))
			groups = append(groups, scanAllBookings(passes, bookings))
		}
	}
	return sortAfterEach(groups)
}

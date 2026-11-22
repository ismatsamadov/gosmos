package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// elevationCalls counts the questions asked through elevation, so a test can see how often a planner asks.
var elevationCalls atomic.Int64

// elevation returns sat's elevation over st as a function of time, counting each question in elevationCalls.
func elevation(sat Satellite, el Elements, st Station) func(time.Time) float64 {
	return func(t time.Time) float64 {
		elevationCalls.Add(1)
		return elevationDeg(sat, el, st, t)
	}
}

// crossing returns the moment between lo and hi when el crosses maskDeg, to a millionth of a second.
// At one of the two, el is below the mask; at the other, at or above it.
func crossing(el func(time.Time) float64, lo, hi time.Time, maskDeg float64) time.Time {
	up := el(hi) >= maskDeg
	for hi.Sub(lo) > time.Microsecond {
		mid := lo.Add(hi.Sub(lo) / 2)
		if (el(mid) >= maskDeg) == up {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
}

// toSecond rounds t to the nearest whole second, and a half second to the even one.
func toSecond(t time.Time) time.Time {
	s := t.Truncate(time.Second)
	if rest := t.Sub(s); rest > time.Second/2 || rest == time.Second/2 && s.Unix()%2 != 0 {
		s = s.Add(time.Second)
	}
	return s
}

// scanEverySecond returns el's passes over maskDeg that start and end from from to to, asking el about
// every second. Where two answers differ, crossing finds the moment: a rise is an AOS, and the fall
// after it the LOS. A pass under 2 minutes, to the second, is too short to book.
func scanEverySecond(el func(time.Time) float64, from, to time.Time, maskDeg float64) []contact {
	var passes []contact
	var aos time.Time
	was := el(from) >= maskDeg
	for t := from.Add(time.Second); t.Before(to); t = t.Add(time.Second) {
		now := el(t) >= maskDeg
		if now == was {
			continue
		}
		at := crossing(el, t.Add(-time.Second), t, maskDeg)
		if now {
			aos = at
		} else if !aos.IsZero() && toSecond(at).Sub(toSecond(aos)) >= 2*time.Minute {
			passes = append(passes, contact{aos, at})
		}
		was = now
	}
	return passes
}

// peak returns the highest elevation el gives at aos and at each whole second after it until los, to a
// tenth of a degree, and the moment of it. The elevation climbs to one top, then falls, so halving the
// seconds by which way it is heading finds the top. Every version of the planner uses peak, so the
// same pass always gets the same MaxElevDeg and the same ID.
func peak(el func(time.Time) float64, aos, los time.Time) (float64, time.Time) {
	after := func(s int) time.Time { return aos.Add(time.Duration(s) * time.Second) }
	lo, hi := 0, int(los.Sub(aos)/time.Second)
	for lo < hi {
		mid := (lo + hi) / 2
		if el(after(mid+1)) > el(after(mid)) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	top, _ := strconv.ParseFloat(strconv.FormatFloat(el(after(lo)), 'f', 1, 64), 64)
	return top, after(lo)
}

// passViews turns the satellite's contacts with the station into passes: AOS and LOS to the nearest
// second, the peak's elevation, and an ID from the orbit number at the peak.
func passViews(sat Satellite, el Elements, st Station, cs []contact) []PassView {
	elev := elevation(sat, el, st)
	views := make([]PassView, 0, len(cs))
	for _, c := range cs {
		top, at := peak(elev, c.AOS, c.LOS)
		views = append(views, PassView{
			ID:   fmt.Sprintf("%s/%s/%05d", sat.ID, st.Code, orbitNumber(sat, el, st, at)),
			Pass: Pass{SatID: sat.ID, Station: st.Code, AOS: toSecond(c.AOS), LOS: toSecond(c.LOS), MaxElevDeg: top},
		})
	}
	return views
}

// scanAllBookings gives each pass its booking, reading through every booking for every pass.
func scanAllBookings(passes []PassView, bookings []Booking) []Planned {
	planned := make([]Planned, 0, len(passes))
	for _, p := range passes {
		id := ""
		for _, b := range bookings {
			if b.SatID == p.SatID && b.Station == p.Station && b.AOS.Equal(p.AOS) {
				id = b.ID
				break
			}
		}
		planned = append(planned, Planned{PassView: p, Booking: id})
	}
	return planned
}

// earlier orders planned passes by AOS, then by ID.
func earlier(a, b Planned) int {
	if c := a.AOS.Compare(b.AOS); c != 0 {
		return c
	}
	return strings.Compare(a.ID, b.ID)
}

// sortAfterEach adds every group's passes to the day one at a time, sorting the whole day after each.
func sortAfterEach(groups [][]Planned) []Planned {
	var day []Planned
	for _, g := range groups {
		for _, p := range g {
			day = append(day, p)
			slices.SortFunc(day, earlier)
		}
	}
	return day
}

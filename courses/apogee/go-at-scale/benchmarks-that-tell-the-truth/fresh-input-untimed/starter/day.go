package main

import (
	"slices"
	"time"
)

// june28 is midnight UTC on Thursday 28 June 2029, take-over day.
var june28 = time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)

// stationDay returns the station with code and its candidates in
// the day from midnight day: the passes of every satellite it can
// serve that climb to 10° or more, satellite by satellite in fleet
// order, so not in AOS order.
func stationDay(code string, day time.Time) (Site, []PassView) {
	n := newNetwork()
	at := slices.IndexFunc(n.Sites, func(s Site) bool {
		return s.Code == code
	})
	s := n.Sites[at]
	var cands []PassView
	for i, e := range n.Elements {
		if !canServe(s, n.Satellites[i]) {
			continue
		}
		for _, p := range predict(e, s, day, day.AddDate(0, 0, 1)) {
			if p.MaxElevDeg >= 10 {
				cands = append(cands, p)
			}
		}
	}
	return s, cands
}

// dayCounts is what planDay has done. The tests read it.
type dayCounts struct {
	calls    int // planDay's calls
	cands    int // candidates handed to it, in all
	sorted   int // calls whose candidates came already sorted
	contacts int // contacts it returned, in all
}

// planned counts planDay's work since the tests last cleared it.
var planned dayCounts

// planDay plans a day's candidates at sites: it sorts cands in
// place, by AOS then pass ID, then gives each site's passes dishes
// with assign.
func planDay(cands []PassView, sites []Site, turnaround time.Duration) []Contact {
	planned.calls++
	planned.cands += len(cands)
	if slices.IsSortedFunc(cands, byAOS) {
		planned.sorted++
	}
	slices.SortFunc(cands, byAOS)
	var out []Contact
	for _, s := range sites {
		here := func(yield func(PassView) bool) {
			for _, p := range cands {
				if p.Station == s.Code && !yield(p) {
					return
				}
			}
		}
		contacts, _ := assign(s, here, nil, turnaround)
		out = append(out, contacts...)
	}
	planned.contacts += len(out)
	return out
}

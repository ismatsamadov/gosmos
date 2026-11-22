package main

import (
	"slices"
	"time"
)

// stationRef plans one station of w the plain way, with no context:
// every candidate predicted at once, less the stowed passes, sorted by
// byAOS and given dishes by assign, dish d+1 free from free[d]. It is
// the tests' answer key.
func stationRef(n *Network, s Site, w Window, free []time.Time, stows []Stow) ([]Contact, int) {
	var cands []PassView
	for i, e := range n.Elements {
		if !canServe(s, n.Satellites[i]) {
			continue
		}
		for _, p := range predict(e, s, w.From, w.To) {
			if p.MaxElevDeg >= 10 && !stowedDuring(p.Pass, stows) {
				cands = append(cands, p)
			}
		}
	}
	slices.SortFunc(cands, byAOS)
	return assign(s, slices.Values(cands), free, 2*time.Minute)
}

// planRef plans every station of w with stationRef, one after another,
// and sorts their contacts by AOS, then pass ID.
func planRef(n *Network, w Window, free map[string][]time.Time, stows []Stow) *Schedule {
	sch := &Schedule{Window: w}
	for _, s := range n.Sites {
		contacts, dropped := stationRef(n, s, w, free[s.Code], stows)
		sch.Contacts = append(sch.Contacts, contacts...)
		sch.Dropped += dropped
	}
	slices.SortFunc(sch.Contacts, func(a, b Contact) int { return byAOS(a.PassView, b.PassView) })
	return sch
}

// handoverRef returns when each dish is free at at: the latest LOS plus
// two minutes of cur's contacts on it with AOS before at.
func handoverRef(cur *Schedule, at time.Time) map[string][]time.Time {
	free := map[string][]time.Time{}
	for _, s := range newNetwork().Sites {
		free[s.Code] = make([]time.Time, s.Dishes)
	}
	for _, c := range cur.Contacts {
		if c.AOS.Before(at) && c.LOS.Add(2*time.Minute).After(free[c.Station][c.Dish-1]) {
			free[c.Station][c.Dish-1] = c.LOS.Add(2 * time.Minute)
		}
	}
	return free
}

// oldPlanner returns the old planner's schedule of a day: the day
// planned from scratch by the same rule, as passplan has planned one
// day at a time since 2023.
func oldPlanner(day time.Time) *Schedule {
	return planRef(newNetwork(), Window{From: day, To: day.AddDate(0, 0, 1)}, nil, nil)
}

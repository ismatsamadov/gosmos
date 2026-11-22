package main

import (
	"cmp"
	"iter"
	"slices"
	"strings"
)

// planOld is the planner Apogee has flown since 2023, on passplan
// 2's pass model and planning rule: it predicts every pair of
// station and satellite into one slice, sorts it by station, then
// AOS, then pass ID, and assigns station by station, growing its
// slices as it goes and keeping each contact's dish in a map keyed
// by pass ID.
func planOld(n *Network, w Window, cfg Config) []Contact {
	var all []PassView
	for _, s := range n.Sites {
		for _, e := range n.Elements {
			all = append(all, predict(e, s, w.From, w.To)...)
		}
	}
	slices.SortFunc(all, func(a, b PassView) int {
		c := strings.Compare(a.Station, b.Station)
		return cmp.Or(c, byAOS(a, b))
	})
	sites := map[string]Site{}
	for _, s := range n.Sites {
		sites[s.Code] = s
	}
	sats := map[string]Satellite{}
	for _, sat := range n.Satellites {
		sats[sat.ID] = sat
	}
	dish := map[string]int{} // pass ID → its dish
	for i := 0; i < len(all); {
		s := sites[all[i].Station]
		var mine []PassView
		for ; i < len(all) && all[i].Station == s.Code; i++ {
			p := all[i]
			high := p.MaxElevDeg >= cfg.MinElevDeg
			if high && canServe(s, sats[p.SatID]) {
				mine = append(mine, p)
			}
		}
		cands := slices.Values(mine)
		contacts, _ := assign(s, cands, nil, cfg.Turnaround)
		for _, c := range contacts {
			dish[c.ID] = c.Dish
		}
	}
	var out []Contact
	for _, p := range all {
		if d, ok := dish[p.ID]; ok {
			out = append(out, Contact{PassView: p, Dish: d})
		}
	}
	slices.SortFunc(out, func(a, b Contact) int {
		return byAOS(a.PassView, b.PassView)
	})
	return out
}

// planNew is passplan 2's planner. Station by station, it gives
// dishes to the candidates stationPasses merges with MergeBy, then
// merges the stations' contacts into one list by AOS, then pass ID.
// It sorts nothing, and makes its slices their full size up front.
func planNew(n *Network, w Window, cfg Config) []Contact {
	stations := make([]iter.Seq[Contact], 0, len(n.Sites))
	total := 0
	for _, s := range n.Sites {
		cands := stationPasses(n, s, w, predict, cfg.MinElevDeg)
		contacts, _ := assign(s, cands, nil, cfg.Turnaround)
		stations = append(stations, slices.Values(contacts))
		total += len(contacts)
	}
	out := make([]Contact, 0, total)
	byContact := func(a, b Contact) int {
		return byAOS(a.PassView, b.PassView)
	}
	for c := range MergeBy(byContact, stations...) {
		out = append(out, c)
	}
	return out
}

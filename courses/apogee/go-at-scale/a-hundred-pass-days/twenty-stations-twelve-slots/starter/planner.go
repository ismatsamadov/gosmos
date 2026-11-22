package main

import (
	"context"
	"slices"
	"time"
)

// planStation plans s's contacts in w: the candidates stationPasses
// yields from passModel, less the passes a stow covers, each given the
// lowest-numbered dish free at its AOS, dish d+1 free from free[d]. It
// returns the contacts and how many candidates found no dish, or
// ctx.Err() once ctx has ended.
func planStation(ctx context.Context, n *Network, s Site, w Window, free []time.Time, stows []Stow, cfg Config) ([]Contact, int, error) {
	cands := func(yield func(PassView) bool) {
		seen := 0
		for p := range stationPasses(n, s, w, passModel, cfg.MinElevDeg) {
			if seen%1000 == 0 && ctx.Err() != nil {
				return // canceled: assign gets no more
			}
			seen++
			if !stowedDuring(p.Pass, stows) && !yield(p) {
				return
			}
		}
	}
	contacts, dropped := assign(s, cands, free, cfg.Turnaround)
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	return contacts, dropped, nil
}

// byContact orders contacts as byAOS orders their passes.
func byContact(a, b Contact) int { return byAOS(a.PassView, b.PassView) }

// Plan plans every station in w, cfg.Slots stations at a time, each
// dish free from free[code], and merges the stations' contacts into one
// schedule by AOS, then pass ID; the caller sets its Version and Made.
// The first error stops the other stations, and Plan returns it.
func Plan(ctx context.Context, n *Network, w Window, free map[string][]time.Time, stows []Stow, cfg Config) (*Schedule, error) {
	lists := make([][]Contact, len(n.Sites)) // each station's contacts, by AOS
	drops := make([]int, len(n.Sites))
	for i, s := range n.Sites { // TODO: plan cfg.Slots stations at a time, and merge their lists
		var err error
		lists[i], drops[i], err = planStation(ctx, n, s, w, free[s.Code], stows, cfg)
		if err != nil {
			return nil, err
		}
	}
	sch := &Schedule{Window: w, Contacts: slices.Concat(lists...)}
	slices.SortFunc(sch.Contacts, byContact)
	for _, d := range drops {
		sch.Dropped += d
	}
	return sch, nil
}

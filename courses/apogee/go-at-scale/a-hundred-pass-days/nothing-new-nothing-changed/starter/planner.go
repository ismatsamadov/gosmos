package main

import (
	"context"
	"iter"
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
	g, _ := newGroup(ctx, cfg.Slots)
	for i, s := range n.Sites {
		g.Go(func(ctx context.Context) error {
			var err error
			lists[i], drops[i], err = planStation(ctx, n, s, w, free[s.Code], stows, cfg)
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	seqs := make([]iter.Seq[Contact], len(lists))
	for i, l := range lists {
		seqs[i] = slices.Values(l)
	}
	sch := &Schedule{Window: w, Contacts: slices.Collect(MergeBy(byContact, seqs...))}
	for _, d := range drops {
		sch.Dropped += d
	}
	return sch, nil
}

// handover returns when each dish is free at the cut-over at: dish d+1
// of a station is free from free[code][d], the latest LOS plus
// turnaround of cur's contacts on it with AOS before at.
func handover(cur *Schedule, at time.Time, turnaround time.Duration) map[string][]time.Time {
	free := map[string][]time.Time{}
	for _, c := range cur.Contacts {
		if !c.AOS.Before(at) {
			break // by AOS: every contact after it rises later still
		}
		dishes := free[c.Station]
		for len(dishes) < c.Dish {
			dishes = append(dishes, time.Time{})
		}
		if t := c.LOS.Add(turnaround); t.After(dishes[c.Dish-1]) {
			dishes[c.Dish-1] = t
		}
		free[c.Station] = dishes
	}
	return free
}

// replanFrom plans cur's window again from at, with stows: cur's
// contacts with AOS before at stay as they are, and the rest is planned
// from the dishes they leave. It returns a new schedule, whose Dropped
// counts the passes from at on.
func replanFrom(ctx context.Context, n *Network, cur *Schedule, at time.Time, stows []Stow, cfg Config) (*Schedule, error) {
	return Plan(ctx, n, cur.Window, nil, stows, cfg) // TODO: keep cur's contacts before at, and plan from at
}

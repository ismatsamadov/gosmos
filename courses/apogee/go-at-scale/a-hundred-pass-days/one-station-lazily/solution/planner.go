package main

import (
	"context"
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

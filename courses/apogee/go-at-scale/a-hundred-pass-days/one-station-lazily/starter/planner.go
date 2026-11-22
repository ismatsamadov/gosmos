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
	return nil, 0, nil
}

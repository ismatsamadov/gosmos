package main

import (
	"context"
	"sync"
)

// Uploads limits how many uploads run at once to each station's dish loader.
type Uploads struct {
	mu        sync.Mutex
	per       int              // uploads at once to one station
	byStation map[string]Slots // each station's slots, made on first use
}

// newUploads returns a limit of perStation uploads at once to each station.
func newUploads(perStation int) *Uploads {
	return &Uploads{per: perStation, byStation: map[string]Slots{}}
}

// Do runs upload while holding one of station's slots, and returns its error, or ctx's if ctx ends while Do waits for a slot.
func (u *Uploads) Do(ctx context.Context, station string, upload func() error) error {
	return upload() // TODO: hold one of station's slots while upload runs
}

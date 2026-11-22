package main

import "context"

// Slots is a semaphore: a channel with room for as many jobs as may run at once.
type Slots chan struct{}

// newSlots returns slots for n jobs at once.
func newSlots(n int) Slots { return make(Slots, n) }

// Acquire waits for a free slot and takes it, or returns ctx's error if ctx ends first.
func (s Slots) Acquire(ctx context.Context) error {
	return nil // TODO: take a slot by sending into s, unless ctx ends first
}

// Release gives back a slot that Acquire took.
func (s Slots) Release() {}

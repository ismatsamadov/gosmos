package main

import (
	"context"
	"time"
)

// waitUntil waits until now() reaches at, unless ctx ends first.
func waitUntil(ctx context.Context, now func() time.Time, at time.Time) error {
	t := time.NewTimer(at.Sub(now()))
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// every calls f with each tick's time, every d, until ctx ends.
func every(ctx context.Context, d time.Duration, f func(at time.Time)) {
	tick := time.NewTicker(d)
	defer tick.Stop()
	for {
		select {
		case at := <-tick.C:
			f(at)
		case <-ctx.Done():
			return
		}
	}
}

// planWithin runs plan, but for no longer than limit.
func planWithin(ctx context.Context, limit time.Duration, plan func(context.Context) (*Schedule, error)) (*Schedule, error) {
	return plan(ctx) // TODO: a context that ends after limit
}

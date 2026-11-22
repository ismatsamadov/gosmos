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
	f(time.Now()) // TODO: a ticker of d, until ctx ends
}

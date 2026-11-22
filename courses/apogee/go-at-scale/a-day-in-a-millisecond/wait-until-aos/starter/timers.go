package main

import (
	"context"
	"time"
)

// waitUntil waits until now() reaches at, unless ctx ends first.
func waitUntil(ctx context.Context, now func() time.Time, at time.Time) error {
	for now().Before(at) {
		if err := ctx.Err(); err != nil {
			return err
		}
		time.Sleep(time.Minute) // TODO: one timer for at.Sub(now())
	}
	return nil
}

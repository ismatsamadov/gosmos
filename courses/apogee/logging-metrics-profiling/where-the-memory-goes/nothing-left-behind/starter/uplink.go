package main

import (
	"context"
	"errors"
	"fmt"
)

// ErrNoWorker is runPass's error when no worker came free before the pass's LOS.
var ErrNoWorker = errors.New("no free worker before LOS")

// runPass runs one pass on a free worker and waits for it until LOS.
func (u *Uplink) runPass(ctx context.Context, p PassView, lines []string) error {
	buf := make([]byte, 0, 4<<20) // the pass's frames, up to 4 MiB
	done := make(chan result, 1)  // room for one: the send never waits for runPass
	go func() {
		u.workers <- struct{}{}        // take one of four workers; wait while all are busy
		defer func() { <-u.workers }() // give the worker back, however the flight ends
		// TODO: take the worker in a select on ctx, then fly with ctx in runPass itself, with no goroutine
		rep, err := u.fly(context.Background(), p, lines, buf) // the radio never hears that the pass is over
		done <- result{rep, err}
	}()
	select {
	case res := <-done:
		return res.err
	case <-ctx.Done(): // the pass's context ends at LOS
		return fmt.Errorf("pass %s: %w", p.ID, ctx.Err())
	}
}

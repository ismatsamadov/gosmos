package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
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
	ctx, cancel := context.WithTimeout(ctx, limit)
	defer cancel()
	s, err := plan(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, fmt.Errorf("plan: took longer than %v: %w", limit, err)
	}
	return s, err
}

// Debouncer calls fire once, quiet after the last of a burst of Notify calls.
type Debouncer struct {
	quiet time.Duration
	fire  func()
	mu    sync.Mutex  // guards timer
	timer *time.Timer // nil until the first Notify
}

// newDebouncer returns a Debouncer that calls fire quiet after a burst ends.
func newDebouncer(quiet time.Duration, fire func()) *Debouncer {
	return &Debouncer{quiet: quiet, fire: fire}
}

// Notify starts the quiet over: fire runs once it has passed.
func (d *Debouncer) Notify() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer == nil {
		d.timer = time.AfterFunc(d.quiet, d.fire)
		return
	}
	d.timer.Reset(d.quiet)
}

// Stop cancels a pending fire.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
}

// flyDay flies each contact at its AOS, in order, until ctx ends.
func flyDay(ctx context.Context, now func() time.Time, contacts []Contact, fly func(Contact)) error {
	for _, c := range contacts {
		fly(c) // TODO: wait for c's AOS with waitUntil first
	}
	return nil
}

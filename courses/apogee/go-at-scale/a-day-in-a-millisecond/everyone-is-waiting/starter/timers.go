package main

import (
	"sync"
	"time"
)

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

package main

import "sync/atomic"

// Counter counts events; it only goes up. Safe for many goroutines.
type Counter struct{ n uint64 }

// Inc adds one to the count.
func (c *Counter) Inc() {
	_ = new(atomic.Uint64) // TODO: keep n in an atomic.Uint64 and add one to it
}

// Add adds n to the count.
func (c *Counter) Add(n uint64) {}

// Value returns the count so far.
func (c *Counter) Value() uint64 { return 0 }

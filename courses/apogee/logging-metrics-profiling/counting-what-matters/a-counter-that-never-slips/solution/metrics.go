package main

import "sync/atomic"

// Counter counts events; it only goes up. Safe for many goroutines.
type Counter struct{ n atomic.Uint64 }

// Inc adds one to the count.
func (c *Counter) Inc() { c.n.Add(1) }

// Add adds n to the count.
func (c *Counter) Add(n uint64) { c.n.Add(n) }

// Value returns the count so far.
func (c *Counter) Value() uint64 { return c.n.Load() }

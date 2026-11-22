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

// Gauge is a reading that goes up and down.
type Gauge struct{ v atomic.Int64 }

// Set makes the reading v.
func (g *Gauge) Set(v int64) { g.v.Store(v) }

// Add moves the reading by d, down when d is negative.
func (g *Gauge) Add(d int64) { g.v.Add(d) }

// Value returns the reading.
func (g *Gauge) Value() int64 { return g.v.Load() }

// Metrics are a station uplinkd's numbers.
type Metrics struct {
	Passes   CounterVec // passes_total: by satellite and outcome
	Failures CounterVec // failures_total: by satellite and reason
	AckWait  *Histogram // ack_wait_seconds
	Running  Gauge      // passes_running
}

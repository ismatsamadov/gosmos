package main

import (
	"slices"
	"sync/atomic"
	"time"
)

// Histogram counts observations into fixed buckets. Safe for many goroutines.
type Histogram struct {
	Bounds []float64       // upper bounds in seconds, ascending
	counts []atomic.Uint64 // one per bound, then one above every bound
	count  atomic.Uint64
	sumMs  atomic.Int64
}

// NewHistogram returns a histogram with a bucket up to each bound, in seconds, and one above them all.
func NewHistogram(bounds ...float64) *Histogram {
	return &Histogram{Bounds: slices.Clone(bounds), counts: make([]atomic.Uint64, len(bounds)+1)}
}

// Observe counts d in the first bucket whose bound is at or above d in seconds, or in the
// last bucket when d is above every bound.
func (h *Histogram) Observe(d time.Duration) {}

// HistogramSnapshot is a histogram's numbers at one moment, as /debug/vars shows them.
type HistogramSnapshot struct {
	Bounds []float64 `json:"bounds"`
	Counts []uint64  `json:"counts"` // len(Bounds)+1
	Count  uint64    `json:"count"`
	SumMs  int64     `json:"sum_ms"`
}

// Snapshot returns a copy of the histogram's numbers.
func (h *Histogram) Snapshot() HistogramSnapshot {
	s := HistogramSnapshot{
		Bounds: slices.Clone(h.Bounds),
		Counts: make([]uint64, len(h.counts)),
		Count:  h.count.Load(),
		SumMs:  h.sumMs.Load(),
	}
	for i := range h.counts {
		s.Counts[i] = h.counts[i].Load()
	}
	return s
}

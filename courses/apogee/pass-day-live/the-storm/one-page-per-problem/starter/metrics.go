package main

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Counter counts events; it only goes up. Safe for many goroutines.
type Counter struct{ n atomic.Uint64 }

// Inc adds one.
func (c *Counter) Inc() { c.n.Add(1) }

// Value returns the count.
func (c *Counter) Value() uint64 { return c.n.Load() }

// CounterVec keeps a counter per pair of label values, such as satellite and outcome.
type CounterVec struct {
	mu sync.RWMutex
	m  map[[2]string]*Counter // made at first use
}

// With returns the counter for a pair of label values: the same one every time.
func (v *CounterVec) With(a, b string) *Counter {
	key := [2]string{a, b}
	v.mu.RLock()
	c, ok := v.m[key]
	v.mu.RUnlock()
	if ok {
		return c
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if c, ok := v.m[key]; ok {
		return c
	}
	if v.m == nil {
		v.m = map[[2]string]*Counter{}
	}
	c = &Counter{}
	v.m[key] = c
	return c
}

// Snapshot returns every count by its two labels: {"brine-3": {"flown": 23}}.
func (v *CounterVec) Snapshot() map[string]map[string]uint64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out := map[string]map[string]uint64{}
	for k, c := range v.m {
		if out[k[0]] == nil {
			out[k[0]] = map[string]uint64{}
		}
		out[k[0]][k[1]] = c.Value()
	}
	return out
}

// Histogram counts observations into fixed buckets. Safe for many goroutines.
type Histogram struct {
	Bounds []float64       // upper bounds in seconds, ascending
	counts []atomic.Uint64 // one per bound, then one above every bound
	count  atomic.Uint64
	sumMs  atomic.Int64
}

// HistogramSnapshot is a histogram's counts at one moment.
type HistogramSnapshot struct {
	Bounds []float64 `json:"bounds"`
	Counts []uint64  `json:"counts"` // len(Bounds)+1
	Count  uint64    `json:"count"`
	SumMs  int64     `json:"sum_ms"`
}

// NewHistogram returns a histogram with the given upper bounds, in seconds.
func NewHistogram(bounds ...float64) *Histogram {
	return &Histogram{Bounds: bounds, counts: make([]atomic.Uint64, len(bounds)+1)}
}

// Observe counts d in the first bucket whose bound is at least d in seconds.
func (h *Histogram) Observe(d time.Duration) {
	h.counts[sort.SearchFloat64s(h.Bounds, d.Seconds())].Add(1)
	h.count.Add(1)
	h.sumMs.Add(d.Milliseconds())
}

// Snapshot returns the histogram's counts now.
func (h *Histogram) Snapshot() HistogramSnapshot {
	s := HistogramSnapshot{Bounds: h.Bounds, Counts: make([]uint64, len(h.counts)), Count: h.count.Load(), SumMs: h.sumMs.Load()}
	for i := range h.counts {
		s.Counts[i] = h.counts[i].Load()
	}
	return s
}

// Metrics are the week's numbers.
type Metrics struct {
	Passes      CounterVec // passes_total: by satellite and outcome
	Failures    CounterVec // failures_total: by satellite and reason
	Commands    CounterVec // commands_total: by satellite, accepted or refused
	Deliveries  CounterVec // deliveries_total: by customer, delivered or undelivered
	Frames      Counter    // frames_total: TLM frames with a good CRC
	BadFrames   Counter    // bad_frames_total: frames of any type with a bad CRC
	BadLines    Counter    // bad_lines_total: TLM lines the strict parser refused
	FilesLanded Counter    // files_landed_total
	AckWait     *Histogram // ack_wait_seconds: from a command sent to its answer
}

// newMetrics returns the week's numbers, all at zero.
func newMetrics() *Metrics {
	return &Metrics{AckWait: NewHistogram(0.5, 1, 2, 5, 10, 30, 60, 120, 300)}
}

// write saves reports/metrics.json in dir: one JSON object, in the shapes the services' metrics use.
func (m *Metrics) write(dir string) error {
	data, err := json.MarshalIndent(struct {
		Passes      map[string]map[string]uint64 `json:"passes_total"`
		Failures    map[string]map[string]uint64 `json:"failures_total"`
		Frames      uint64                       `json:"frames_total"`
		BadFrames   uint64                       `json:"bad_frames_total"`
		BadLines    uint64                       `json:"bad_lines_total"`
		Commands    map[string]map[string]uint64 `json:"commands_total"`
		FilesLanded uint64                       `json:"files_landed_total"`
		Deliveries  map[string]map[string]uint64 `json:"deliveries_total"`
		AckWait     HistogramSnapshot            `json:"ack_wait_seconds"`
	}{m.Passes.Snapshot(), m.Failures.Snapshot(), m.Frames.Value(), m.BadFrames.Value(), m.BadLines.Value(),
		m.Commands.Snapshot(), m.FilesLanded.Value(), m.Deliveries.Snapshot(), m.AckWait.Snapshot()}, "", "  ")
	if err != nil {
		return err
	}
	return saveWhole(filepath.Join(dir, "reports", "metrics.json"), append(data, '\n'))
}

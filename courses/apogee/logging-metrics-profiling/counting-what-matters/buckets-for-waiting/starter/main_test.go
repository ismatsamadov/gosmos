package main

import (
	"slices"
	"sync"
	"testing"
	"time"
)

// bounds are ack_wait_seconds' buckets, in seconds.
var bounds = []float64{0.5, 1, 2, 5, 10, 30, 60, 120, 300}

func TestObserveBuckets(t *testing.T) {
	tests := []struct {
		name   string
		d      time.Duration
		bucket int
	}{
		{"an answer at once", 0, 0},
		{"exactly half a second", 500 * time.Millisecond, 0},
		{"just over half a second", 501 * time.Millisecond, 1},
		{"exactly 2 seconds", 2 * time.Second, 2},
		{"Brine-3's 71 seconds after the restart on 8 November", 71 * time.Second, 7},
		{"exactly 5 minutes", 300 * time.Second, 8},
		{"301 seconds, above every bound", 301 * time.Second, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHistogram(bounds...)
			h.Observe(tt.d)
			want := make([]uint64, len(bounds)+1)
			want[tt.bucket] = 1
			if got := h.Snapshot().Counts; !slices.Equal(got, want) {
				t.Errorf("counts after Observe(%v) = %v, want %v", tt.d, got, want)
			}
		})
	}
}

func TestObserveCountAndSum(t *testing.T) {
	h := NewHistogram(bounds...)
	h.Observe(71 * time.Second)
	h.Observe(1500 * time.Millisecond)
	h.Observe(999 * time.Microsecond)
	s := h.Snapshot()
	if s.Count != 3 {
		t.Errorf("count after Observe(1m11s), Observe(1.5s) and Observe(999µs) = %d, want 3", s.Count)
	}
	if s.SumMs != 72500 {
		t.Errorf("sum_ms after Observe(1m11s), Observe(1.5s) and Observe(999µs) = %d, want 72500: each wait's whole milliseconds, its fraction dropped", s.SumMs)
	}
}

// TestObserveFromManyGoroutines has 30 goroutines observe 100 waits each at once, ten in
// every bucket.
func TestObserveFromManyGoroutines(t *testing.T) {
	waits := []time.Duration{ // one for each bucket
		300 * time.Millisecond, 800 * time.Millisecond, 1500 * time.Millisecond, 3 * time.Second,
		7 * time.Second, 20 * time.Second, 45 * time.Second, 90 * time.Second,
		200 * time.Second, 400 * time.Second,
	}
	h := NewHistogram(bounds...)
	var wg sync.WaitGroup
	for range 30 {
		wg.Go(func() {
			for i := range 100 {
				h.Observe(waits[i%len(waits)])
			}
		})
	}
	wg.Wait()
	s := h.Snapshot()
	want := []uint64{300, 300, 300, 300, 300, 300, 300, 300, 300, 300}
	if !slices.Equal(s.Counts, want) {
		t.Errorf("counts after 30 goroutines × 100 Observe = %v, want %v", s.Counts, want)
	}
	if s.Count != 3000 {
		t.Errorf("count after 30 goroutines × 100 Observe = %d, want 3000", s.Count)
	}
	if s.SumMs != 230280000 {
		t.Errorf("sum_ms after 30 goroutines × 100 Observe = %d, want 230280000", s.SumMs)
	}
}

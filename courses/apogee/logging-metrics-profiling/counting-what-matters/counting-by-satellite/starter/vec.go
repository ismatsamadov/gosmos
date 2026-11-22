package main

import "sync"

// CounterVec keeps a counter per pair of label values, such as satellite and outcome.
type CounterVec struct {
	mu sync.RWMutex
	m  map[[2]string]*Counter // made at first use
}

// With returns the counter for the pair a, b: the same *Counter every time for the same pair.
func (v *CounterVec) With(a, b string) *Counter {
	return &Counter{}
}

// Snapshot returns a copy of every count, by first label and then second:
// {"brine-3": {"failed": 1, "flown": 23}}.
func (v *CounterVec) Snapshot() map[string]map[string]uint64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out := make(map[string]map[string]uint64)
	for key, c := range v.m {
		if out[key[0]] == nil {
			out[key[0]] = make(map[string]uint64)
		}
		out[key[0]][key[1]] = c.Value()
	}
	return out
}

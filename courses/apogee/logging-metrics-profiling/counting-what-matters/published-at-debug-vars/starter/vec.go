package main

import "sync"

// CounterVec keeps a counter per pair of label values, such as satellite and outcome.
type CounterVec struct {
	mu sync.RWMutex
	m  map[[2]string]*Counter // made at first use
}

// With returns the counter for the pair a, b: the same *Counter every time for the same pair.
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
		v.m = make(map[[2]string]*Counter)
	}
	v.m[key] = &Counter{}
	return v.m[key]
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

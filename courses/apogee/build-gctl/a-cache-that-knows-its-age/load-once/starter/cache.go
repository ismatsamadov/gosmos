package main

import "time"

// passCache keeps the night's passes after loading them, for a time to live.
type passCache struct {
	load     func() ([]Pass, error) // reads the night's passes
	now      func() time.Time       // the clock: time.Now, or a test's own
	ttl      time.Duration          // how long loaded passes stay fresh
	cached   []Pass
	loadedAt time.Time
	loaded   bool
}

// passes returns the night's passes, calling load only when none are loaded yet.
func (c *passCache) passes() ([]Pass, error) {
	return nil, nil
}

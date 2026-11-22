package main

import (
	"slices"
	"sync"
	"time"
)

// passCache keeps the night's passes after loading them, for a time to live. Its methods are
// safe to call from several goroutines at once.
type passCache struct {
	mu       sync.Mutex
	load     func() ([]Pass, error) // reads the night's passes
	now      func() time.Time       // the clock: time.Now, or a test's own
	ttl      time.Duration          // how long loaded passes stay fresh
	cached   []Pass
	loadedAt time.Time
	loaded   bool
}

// passes returns a copy of the night's passes, loading them when none are loaded or they are
// stale.
func (c *passCache) passes() ([]Pass, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.loaded || !c.now().Before(c.loadedAt.Add(c.ttl)) {
		if err := c.reload(); err != nil {
			return nil, err
		}
	}
	return slices.Clone(c.cached), nil
}

// refresh loads the night's passes now, keeps them and records when.
func (c *passCache) refresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reload()
}

// reload loads the night's passes, keeps them and records when; the caller holds c.mu.
func (c *passCache) reload() error {
	ps, err := c.load()
	if err != nil {
		return err
	}
	c.cached, c.loadedAt, c.loaded = ps, c.now(), true
	return nil
}

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

// passes returns the night's passes, loading them when none are loaded or they are stale.
func (c *passCache) passes() ([]Pass, error) {
	if !c.loaded || !c.now().Before(c.loadedAt.Add(c.ttl)) {
		if err := c.refresh(); err != nil {
			return nil, err
		}
	}
	return c.cached, nil
}

// refresh loads the night's passes now, keeps them and records when.
func (c *passCache) refresh() error {
	ps, err := c.load()
	if err != nil {
		return err
	}
	c.cached, c.loadedAt, c.loaded = ps, c.now(), true
	return nil
}

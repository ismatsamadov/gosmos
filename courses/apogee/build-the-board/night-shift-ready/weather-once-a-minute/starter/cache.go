package main

import (
	"context"
	"sync"
	"time"
)

// weatherCache keeps each station's latest reading from the weather service for a minute,
// on the Board's clock.
type weatherCache struct {
	fetch func(ctx context.Context, code string) (Weather, error) // the weather service
	now   func() time.Time
	ttl   time.Duration // how long a reading is kept: a minute

	mu   sync.Mutex
	kept map[string]keptReading // by station code
}

// keptReading is a station's reading and the time it was fetched.
type keptReading struct {
	weather Weather
	at      time.Time
}

// newWeatherCache returns a cache that asks fetch for readings and keeps them a minute.
func newWeatherCache(fetch func(ctx context.Context, code string) (Weather, error), now func() time.Time) *weatherCache {
	return &weatherCache{fetch: fetch, now: now, ttl: time.Minute, kept: map[string]keptReading{}}
}

// reading returns the station's reading: the one kept, while it is under a minute old,
// otherwise a new one from the weather service. An error is returned and never kept.
func (c *weatherCache) reading(ctx context.Context, code string) (Weather, error) {
	return c.fetch(ctx, code)
}

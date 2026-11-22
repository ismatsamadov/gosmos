package main

import (
	"fmt"
	"maps"
	"sync"
)

// StatusCache holds every station's latest status, for many readers and one writer.
type StatusCache struct {
	mu     sync.RWMutex
	status map[string]string
}

// NewStatusCache returns an empty cache.
func NewStatusCache() *StatusCache {
	return &StatusCache{status: map[string]string{}}
}

// Set records a station's status.
func (c *StatusCache) Set(code, status string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status[code] = status
}

// Get returns a station's status, or "" if it has none.
func (c *StatusCache) Get(code string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status[code]
}

// Snapshot returns a copy of every station's status, for the Board.
func (c *StatusCache) Snapshot() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return maps.Clone(c.status)
}

func main() {
	c := NewStatusCache()
	c.Set("YST", "stowed: wind 81 km/h")
	c.Set("SQT", "tracking furrow-a")
	fmt.Println(c.Snapshot())
}

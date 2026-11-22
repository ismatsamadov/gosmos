package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"sync"
	"time"
)

// Quota counts each customer's bookings in the current calendar month against its plan, and keeps the count in a file, so a restart forgets nothing.
type Quota struct {
	plans map[string]int // bookings a month, by customer ID; Apogee's own have no plan
	now   func() time.Time
	path  string // usage.json in the data folder
	mu    sync.Mutex
	usage usage
}

// openQuota opens the count saved at path, on the clock now: a missing file is a new count, a broken one an error.
func openQuota(path string, plans map[string]int, now func() time.Time) (*Quota, error) {
	q := &Quota{plans: plans, now: now, path: path, usage: usage{Used: map[string]int{}}}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return q, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &q.usage); err != nil {
		return nil, fmt.Errorf("usage %s: %w", path, err)
	}
	if q.usage.Used == nil {
		q.usage.Used = map[string]int{}
	}
	return q, nil
}

// Take counts one booking for customer and saves the count, or refuses the booking when its plan for this month is used up. A customer without a plan is never counted.
func (q *Quota) Take(customer string) error {
	plan, ok := q.plans[customer]
	if !ok {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	before := usage{Month: q.usage.Month, Used: maps.Clone(q.usage.Used)}
	if err := q.usage.take(customer, plan, q.now()); err != nil {
		return err
	}
	if err := q.save(); err != nil {
		q.usage = before
		return err
	}
	return nil
}

// save writes the whole count to its file through replaceFile. The caller holds q.mu.
func (q *Quota) save() error {
	data, err := json.Marshal(q.usage)
	if err != nil {
		return err
	}
	return replaceFile(q.path, data)
}

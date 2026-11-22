package main

import "sync/atomic"

// Stats counts what passplan 2's planning runs did, added from any number of goroutines at once.
type Stats struct {
	Planned atomic.Int64 // contacts planned
	Dropped atomic.Int64 // passes no dish could take
	Runs    atomic.Int64 // runs counted
}

// Add counts one run that planned planned contacts and dropped dropped passes.
func (s *Stats) Add(planned, dropped int) {
	s.Planned.Add(int64(planned))
	s.Dropped.Add(int64(dropped))
	s.Runs.Add(1)
}

// Snapshot returns the counts so far.
func (s *Stats) Snapshot() (planned, dropped, runs int64) {
	return s.Planned.Load(), s.Dropped.Load(), s.Runs.Load()
}

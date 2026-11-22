package main

import (
	"errors"
	"slices"
	"time"
)

// ErrStowed means a station's dish is stowed: none of its remaining passes today can run.
var ErrStowed = errors.New("dish stowed")

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID   string
	Station string    // station code
	AOS     time.Time // acquisition of signal
	LOS     time.Time // loss of signal
}

// Report counts how a day's passes went.
type Report struct {
	Done, Failed, Skipped int
}

// Add adds another report's counts to r.
func (r *Report) Add(other Report) {
	r.Done += other.Done
	r.Failed += other.Failed
	r.Skipped += other.Skipped
}

// byStation groups passes by station, each station's passes in AOS order.
func byStation(passes []Pass) map[string][]Pass {
	days := map[string][]Pass{}
	for _, p := range passes {
		days[p.Station] = append(days[p.Station], p)
	}
	for _, day := range days {
		slices.SortFunc(day, func(a, b Pass) int { return a.AOS.Compare(b.AOS) })
	}
	return days
}

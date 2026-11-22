package main

import (
	"sync/atomic"
	"time"
)

// inputs are what planDay plans from.
type inputs struct {
	sats     []Satellite
	els      map[string]Elements
	stations []Station
	bookings []Booking
	from, to time.Time
}

// setups counts benchInputs's calls, so a test can see how often a benchmark builds its inputs.
var setups atomic.Int64

// benchInputs builds what Nkem measures: Brine-3 over Cape Wren for the three hours from 05:10 on
// Tuesday 16 January.
func benchInputs() inputs {
	setups.Add(1)
	from := time.Date(2029, 1, 16, 5, 10, 0, 0, time.UTC)
	return inputs{
		sats:     []Satellite{satellite("brine-3")},
		els:      elements(),
		stations: []Station{station("CWR")},
		from:     from,
		to:       from.Add(3 * time.Hour),
	}
}

package main

import "time"

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID   string
	Station string    // station code
	AOS     time.Time // acquisition of signal
	LOS     time.Time // loss of signal
}

package main

import "time"

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	// TODO: tag each field with its key in the export.
	ID      string
	SatID   string
	Station string
	AOS     time.Time // when the satellite rises into view
	LOS     time.Time // when it sets
}

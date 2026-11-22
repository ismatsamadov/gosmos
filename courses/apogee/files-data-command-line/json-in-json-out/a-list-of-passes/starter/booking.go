package main

import "time"

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string    `json:"id"`
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"` // when the satellite rises into view
	LOS     time.Time `json:"los"` // when it sets
}

package main

import "time"

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string
	Station string
	AOS     time.Time // when the satellite rises into view
	LOS     time.Time // when it sets
}

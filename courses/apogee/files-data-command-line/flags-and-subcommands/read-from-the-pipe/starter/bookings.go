package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string    `json:"id"`
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
}

// readBookings decodes a booking log, one JSON booking per line, returning those read before any error.
func readBookings(r io.Reader) ([]Booking, error) {
	dec := json.NewDecoder(r)
	var bookings []Booking
	for {
		var b Booking
		err := dec.Decode(&b)
		if err == io.EOF {
			return bookings, nil
		}
		if err != nil {
			return bookings, fmt.Errorf("booking %d: %w", len(bookings)+1, err)
		}
		bookings = append(bookings, b)
	}
}

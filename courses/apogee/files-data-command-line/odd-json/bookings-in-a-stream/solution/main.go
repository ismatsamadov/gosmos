package main

import (
	"encoding/json"
	"io"
	"time"
)

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string    `json:"id"`
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"` // when the satellite rises into view
	LOS     time.Time `json:"los"` // when it sets
	Note    string    `json:"note,omitempty"`
}

// decodeStream decodes bookings from r, one JSON value after another, until the stream ends.
func decodeStream(r io.Reader) ([]Booking, error) {
	dec := json.NewDecoder(r)
	var bookings []Booking
	for {
		var b Booking
		err := dec.Decode(&b)
		if err == io.EOF {
			return bookings, nil
		}
		if err != nil {
			return bookings, err
		}
		bookings = append(bookings, b)
	}
}

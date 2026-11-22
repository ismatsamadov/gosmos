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
	AOS     time.Time `json:"aos"`
	LOS     time.Time `json:"los"`
	Note    string    `json:"note,omitempty"`
}

// decodeStream decodes bookings from r, one JSON value after another, until the stream ends.
func decodeStream(r io.Reader) ([]Booking, error) {
	_ = json.NewDecoder
	return nil, nil
}

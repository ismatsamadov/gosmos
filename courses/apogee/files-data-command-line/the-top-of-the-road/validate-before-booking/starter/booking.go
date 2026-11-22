package main

import "time"

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string    `json:"id"`
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"` // when the satellite rises into view
	LOS     time.Time `json:"los"` // when it sets
	Note    string    `json:"note,omitempty"`
}

// Band is a radio band.
type Band string

// Station is one of Apogee's ground stations.
type Station struct {
	Code    string // "YST": three capital letters
	Name    string // "Ystmark"
	MaskDeg float64
	Bands   []Band
}

package main

import "time"

// Band is a radio band.
type Band string

const (
	UHF   Band = "UHF"
	SBand Band = "S"
	XBand Band = "X"
)

// Station is one of Apogee's ground stations.
type Station struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	LatDeg     float64 `json:"lat_deg"`
	LonDeg     float64 `json:"lon_deg"`
	MaskDeg    float64 `json:"mask_deg"`
	Bands      []Band  `json:"bands"`
	Stowed     bool    `json:"stowed"`
	StowReason string  `json:"stow_reason,omitempty"`
}

// Orbit is a circular orbit.
type Orbit struct {
	AltKm  float64
	IncDeg float64
}

// Mode is a satellite's operating mode.
type Mode string

// Satellite is a spacecraft Apogee talks to.
type Satellite struct {
	ID       string
	Name     string
	Operator string
	Orbit    Orbit
	Bands    []Band
	Mode     Mode
}

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID      string    `json:"sat_id"`
	Station    string    `json:"station"`
	AOS        time.Time `json:"aos"`
	LOS        time.Time `json:"los"`
	MaxElevDeg float64   `json:"max_elev_deg"`
}

// PassView is a predicted pass with its pass ID.
type PassView struct {
	ID string `json:"id"`
	Pass
}

// TelemetryFrame is one decoded housekeeping frame.
type TelemetryFrame struct {
	SatID string
	Seq   uint32
	At    time.Time
	BattV float64
	TempC float64
	Mode  Mode
}

// Booking is one reserved pass.
type Booking struct {
	ID      string    `json:"id"` // "bk-4501", given by the service
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
	LOS     time.Time `json:"los"`
	Note    string    `json:"note,omitempty"`
}

// Landing is dlinkd's word that a booking's file has been delivered.
type Landing struct {
	BookingID string    // "bk-4912"
	File      string    // "furrow-a-20280201-1402.tif": satellite, AOS day and time
	Bytes     int64     // the file's size
	At        time.Time // when the file was complete, UTC
}

// Event is what the service tells a customer when something happens to a booking.
type Event struct {
	ID        string    `json:"id"`   // "evt-000118": the same on every attempt
	Type      string    `json:"type"` // "data.landed", the only type so far
	BookingID string    `json:"booking_id"`
	SatID     string    `json:"sat_id"`
	Station   string    `json:"station"`
	File      string    `json:"file"`
	Bytes     int64     `json:"bytes"`
	LandedAt  time.Time `json:"landed_at"`
}

// Window is the span a plan covers: From included, To excluded.
type Window struct {
	From time.Time
	To   time.Time
}

// Stow is a span when a station's dishes are parked. A zero To lasts until further notice.
type Stow struct {
	Station string    // "CWR"
	From    time.Time // 2029-06-28T09:47:00Z
	To      time.Time // 2029-06-28T13:23:00Z
	Reason  string    // "wind 74 km/h"
}

// Contact is a pass given a dish: what a station will fly.
type Contact struct {
	PassView     // the pass and its ID
	Dish     int `json:"dish"` // 1 or 2
}

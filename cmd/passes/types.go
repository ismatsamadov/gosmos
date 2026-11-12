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

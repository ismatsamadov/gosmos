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

// Site is a station as passplan sees it.
type Site struct {
	Station
	Dishes int       // the dishes passplan can give passes to
	Live   time.Time // the first moment it flies a pass; zero means always
}

// Elements place a satellite on its orbit for the pass model.
type Elements struct {
	SatID    string
	Orbit    Orbit     // altitude and inclination
	Epoch    time.Time // the moment the numbers below describe; no pass before it
	Rev0     int       // the orbit number at Epoch
	NodeDeg  float64   // east longitude of the northbound equator crossing, at Epoch
	PhaseDeg float64   // how far past that crossing the satellite is at Epoch, 0 to 360
	// RevsPerDay is the satellite's orbits a day, where tracking measured them; zero means
	// the period its altitude gives.
	RevsPerDay float64
}

// Network is everything passplan plans for.
type Network struct {
	Sites      []Site      // the twenty, in table order
	Satellites []Satellite // the eighty, in fleet order
	Elements   []Elements  // the same eighty, in the same order
}

// Window is the span a plan covers: From included, To excluded.
type Window struct {
	From time.Time
	To   time.Time
}

// Contact is a pass given a dish: what a station will fly.
type Contact struct {
	PassView     // the pass and its ID
	Dish     int `json:"dish"` // 1 or 2
}

// Stow is a span when a station's dishes are parked. A zero To lasts until further notice.
type Stow struct {
	Station string    // "CWR"
	From    time.Time // 2029-06-28T09:47:00Z
	To      time.Time // 2029-06-28T13:23:00Z
	Reason  string    // "wind 74 km/h"
}

// Schedule is one published plan. Nobody changes it after it is published.
type Schedule struct {
	Version  uint64    // the run that made it
	Window             // what it plans
	Made     time.Time // when its run started
	Contacts []Contact // every contact with AOS in the window, by AOS then ID
	Dropped  int       // candidates no dish could take
}

// Config is passplan 2's settings.
type Config struct {
	Slots      int           `check:"min=1,max=64" default:"12"`
	Days       int           `check:"min=1,max=120" default:"100"`
	MinElevDeg float64       `check:"min=0,max=45" default:"10"`
	Turnaround time.Duration `check:"min=0s,max=10m" default:"2m"`
	Replan     time.Duration `check:"min=1m,max=1h" default:"10m"`
	RunLimit   time.Duration `check:"min=1s,max=5m" default:"30s"`
	Queue      int           `check:"min=1,max=4096" default:"256"`
	Quiet      time.Duration `check:"min=0s,max=1m" default:"5s"`
}

// Event is news for the replanner: a tick, or a stow that begins or ends.
type Event struct {
	At   time.Time
	Stow *Stow // nil for a tick; a stow with a zero To begins, one with a To ends it
}

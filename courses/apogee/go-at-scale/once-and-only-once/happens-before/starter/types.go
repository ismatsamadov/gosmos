package main

import "time"

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

// Schedule is one published plan. Nobody changes it after it is published.
type Schedule struct {
	Version  uint64    // the run that made it
	Window             // what it plans
	Made     time.Time // when its run started
	Contacts []Contact // every contact with AOS in the window, by AOS then ID
	Dropped  int       // candidates no dish could take
}

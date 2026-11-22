package main

import "time"

// Station is one of Apogee's ground stations.
type Station struct {
	Code    string  // "YST": three capital letters
	LatDeg  float64 // north is positive
	LonDeg  float64 // east is positive
	MaskDeg float64 // lowest usable elevation
}

// Orbit is a circular orbit.
type Orbit struct {
	AltKm  float64
	IncDeg float64
}

// Satellite is a spacecraft Apogee talks to.
type Satellite struct {
	ID    string // "pipit-1"
	Orbit Orbit
}

// Elements place a satellite on its orbit for the pass model.
type Elements struct {
	Epoch    time.Time // the moment the numbers below describe; no pass before it
	Rev0     int       // the orbit number at Epoch
	NodeDeg  float64   // east longitude of the northbound equator crossing, at Epoch
	PhaseDeg float64   // how far past that crossing the satellite is at Epoch, 0 to 360
	// RevsPerDay is the satellite's orbits a day, where tracking measured them; zero means
	// the period its altitude gives.
	RevsPerDay float64
}

// epoch is when the elements hold: Tuesday 16 January 2029, 00:00 UTC.
var epoch = time.Date(2029, 1, 16, 0, 0, 0, 0, time.UTC)

// fleet returns the seven satellites passplan plans for.
func fleet() []Satellite {
	return []Satellite{
		{ID: "pipit-1", Orbit: Orbit{AltKm: 525, IncDeg: 97.5}},
		{ID: "pipit-2", Orbit: Orbit{AltKm: 550, IncDeg: 97.6}},
		{ID: "pipit-3", Orbit: Orbit{AltKm: 510, IncDeg: 97.4}},
		{ID: "furrow-a", Orbit: Orbit{AltKm: 510, IncDeg: 97.4}},
		{ID: "furrow-b", Orbit: Orbit{AltKm: 510, IncDeg: 97.4}},
		{ID: "brine-3", Orbit: Orbit{AltKm: 480, IncDeg: 53.0}},
		{ID: "cirrine-1", Orbit: Orbit{AltKm: 600, IncDeg: 97.8}},
	}
}

// elements returns every satellite's elements at the epoch, by satellite ID.
func elements() map[string]Elements {
	return map[string]Elements{
		"pipit-1":   {Epoch: epoch, Rev0: 14321, NodeDeg: 201.4926001313, PhaseDeg: 322.4045275738},
		"pipit-2":   {Epoch: epoch, Rev0: 19979, NodeDeg: 234.8439457617, PhaseDeg: 41.6844523733},
		"pipit-3":   {Epoch: epoch, Rev0: 9917, NodeDeg: 149.9999975496, PhaseDeg: 93.2834560607},
		"furrow-a":  {Epoch: epoch, Rev0: 22046, NodeDeg: 294.5055026568, PhaseDeg: 109.4557673805, RevsPerDay: 867.0 / 57},
		"furrow-b":  {Epoch: epoch, Rev0: 22043, NodeDeg: 297.5851207259, PhaseDeg: 287.9312984850, RevsPerDay: 867.0 / 57},
		"brine-3":   {Epoch: epoch, Rev0: 15439, NodeDeg: 94.5696529379, PhaseDeg: 61.6707270863},
		"cirrine-1": {Epoch: epoch, Rev0: 17430, NodeDeg: 282.5236519320, PhaseDeg: 176.5142936516},
	}
}

// allStations returns Apogee's five ground stations, with where they are and their masks.
func allStations() []Station {
	return []Station{
		{Code: "YST", LatDeg: 77.8, LonDeg: 2.0, MaskDeg: 5},
		{Code: "PEL", LatDeg: 46.9, LonDeg: -38.1, MaskDeg: 5},
		{Code: "SQT", LatDeg: -23.4, LonDeg: -104.6, MaskDeg: 3},
		{Code: "LAN", LatDeg: -1.8, LonDeg: 82.5, MaskDeg: 8},
		{Code: "CWR", LatDeg: -53.9, LonDeg: -95.0, MaskDeg: 5},
	}
}

// satellite returns the fleet's satellite with the given ID, or a zero Satellite.
func satellite(id string) Satellite {
	for _, s := range fleet() {
		if s.ID == id {
			return s
		}
	}
	return Satellite{}
}

// station returns the station with the given code, or a zero Station.
func station(code string) Station {
	for _, s := range allStations() {
		if s.Code == code {
			return s
		}
	}
	return Station{}
}

package main

import "fmt"

// Orbit is a circular orbit.
type Orbit struct {
	AltKm  float64
	IncDeg float64
}

// Satellite is a spacecraft Apogee talks to.
type Satellite struct {
	ID    string
	Name  string
	Orbit Orbit
}

// String returns the satellite's name and altitude, the way ops reads them.
func (s Satellite) String() string {
	return fmt.Sprintf("%s (%.0f km)", s.Name, s.Orbit.AltKm)
}

func main() {
	sats := []Satellite{
		{ID: "pipit-1", Name: "Pipit-1", Orbit: Orbit{AltKm: 525, IncDeg: 97.5}},
		{ID: "thistle", Name: "Thistle", Orbit: Orbit{AltKm: 415, IncDeg: 51.6}},
		{ID: "cirrine-1", Name: "Cirrine-1", Orbit: Orbit{AltKm: 600, IncDeg: 97.8}},
	}
	for _, s := range sats {
		fmt.Println(s)
	}
}

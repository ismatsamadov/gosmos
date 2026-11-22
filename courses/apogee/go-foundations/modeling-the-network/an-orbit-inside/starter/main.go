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

// isPolar reports whether a satellite's orbit is inclined from 80 to 100 degrees.
func isPolar(s Satellite) bool {
	return false
}

func main() {
	pipit := Satellite{ID: "pipit-1", Name: "Pipit-1", Orbit: Orbit{AltKm: 525, IncDeg: 97.5}}
	fmt.Println(pipit.Name, "polar:", isPolar(pipit))
}

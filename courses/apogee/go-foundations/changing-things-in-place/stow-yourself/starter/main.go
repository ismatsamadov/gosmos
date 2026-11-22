package main

import "fmt"

// Station is one of Apogee's ground stations.
type Station struct {
	Code       string
	Name       string
	Stowed     bool
	StowReason string
}

// Stow parks the station's dishes and records why.
func (s Station) Stow(reason string) {
	s.Stowed = true
	s.StowReason = reason
}

// Unstow puts the station's dishes back to work.
func (s Station) Unstow() {
	s.Stowed = false
	s.StowReason = ""
}

func main() {
	yst := Station{Code: "YST", Name: "Ystmark"}
	yst.Stow("wind 81 km/h")
	fmt.Printf("%s stowed: %t (%s)\n", yst.Code, yst.Stowed, yst.StowReason)
}

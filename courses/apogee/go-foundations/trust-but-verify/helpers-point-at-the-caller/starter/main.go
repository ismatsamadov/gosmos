package main

import "fmt"

// Earth and physics.
const (
	earthRadiusKm = 6371
	muKm3PerS2    = 398600.4418 // Earth's gravitational parameter
	speedOfLight  = 299792458.0 // meters per second
)

// dopplerHz returns the shift of a signal at freqHz when the range changes at rangeRateMps.
func dopplerHz(freqHz, rangeRateMps float64) float64 {
	return -freqHz * rangeRateMps / speedOfLight
}

func main() {
	fmt.Printf("Pipit-1 rising: %+.2f Hz\n", dopplerHz(437.25e6, -7000))
}

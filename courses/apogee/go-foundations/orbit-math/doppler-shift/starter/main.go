package main

import "fmt"

// speedOfLight is in meters per second.
const speedOfLight = 299792458.0

// dopplerHz returns the shift of a signal at freqHz when the range changes at rangeRateMps.
func dopplerHz(freqHz, rangeRateMps float64) float64 {
	return 0
}

func main() {
	fmt.Printf("Pipit-1 rising: %.2f Hz\n", dopplerHz(437.25e6, -7000))
}

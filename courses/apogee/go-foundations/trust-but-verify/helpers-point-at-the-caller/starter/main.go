package main

import "fmt"

// speedOfLight is in meters per second.
const speedOfLight = 299792458.0

// dopplerHz returns the Doppler shift of a signal at freqHz seen at a range rate in m/s.
func dopplerHz(freqHz, rangeRateMps float64) float64 {
	return -freqHz * rangeRateMps / speedOfLight
}

func main() {
	fmt.Printf("Pipit-1 rising: %+.2f Hz\n", dopplerHz(437.25e6, -7000))
}

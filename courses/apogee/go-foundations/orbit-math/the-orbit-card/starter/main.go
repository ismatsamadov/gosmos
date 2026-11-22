package main

import (
	"fmt"
	"math"
)

// Earth and physics.
const (
	earthRadiusKm = 6371
	muKm3PerS2    = 398600.4418 // Earth's gravitational parameter
	speedOfLight  = 299792458.0 // meters per second
)

// Pipit-2's pass over Salar Quieto.
const (
	pipit2AltKm       = 550.0
	pipit2DownlinkMHz = 2245.0
	passElevDeg       = 10.0
	passRangeRateMps  = -6500.0
)

// radians converts degrees to radians.
func radians(deg float64) float64 {
	return deg * math.Pi / 180
}

// periodMinutes returns the period of a circular orbit at altKm.
func periodMinutes(altKm float64) float64 {
	a := earthRadiusKm + altKm
	return 2 * math.Pi * math.Sqrt(a*a*a/muKm3PerS2) / 60
}

// slantRangeKm returns the distance to a satellite at altKm, seen at elevDeg.
func slantRangeKm(altKm, elevDeg float64) float64 {
	e := radians(elevDeg)
	a := earthRadiusKm + altKm
	b := earthRadiusKm * math.Cos(e)
	return math.Sqrt(a*a-b*b) - earthRadiusKm*math.Sin(e)
}

// dopplerHz returns the shift of a signal at freqHz when the range changes at rangeRateMps.
func dopplerHz(freqHz, rangeRateMps float64) float64 {
	return -freqHz * rangeRateMps / speedOfLight
}

// pathLossDB returns the free-space path loss at freqMHz to a satellite at altKm, seen at elevDeg.
func pathLossDB(altKm, elevDeg, freqMHz float64) float64 {
	d := slantRangeKm(altKm, elevDeg)
	return 20*math.Log10(d) + 20*math.Log10(freqMHz) + 32.44
}

func main() {
	fmt.Println("Pipit-2 over Salar Quieto")
	// TODO: print the period, slant range, Doppler shift and path loss.
}

// Package orbit computes facts about circular orbits.
package orbit

import "math"

const (
	earthRadiusKm = 6371
	muKm3PerS2    = 398600.4418 // Earth's gravitational parameter
)

// cube returns x to the power of three.
func cube(x float64) float64 {
	return x * x * x
}

// periodMinutes returns the period of a circular orbit altKm above the Earth, in minutes.
func periodMinutes(altKm float64) float64 {
	return 2 * math.Pi * math.Sqrt(cube(earthRadiusKm+altKm)/muKm3PerS2) / 60
}

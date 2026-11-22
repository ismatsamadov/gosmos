// Package orbit computes facts about circular orbits.
package orbit

import "math"

const (
	earthRadiusKm = 6371
	muKm3PerS2    = 398600.4418 // Earth's gravitational parameter
)

// PeriodMinutes returns the period of a circular orbit altKm above the Earth, in minutes.
func PeriodMinutes(altKm float64) float64 {
	a := earthRadiusKm + altKm
	return 2 * math.Pi * math.Sqrt(a*a*a/muKm3PerS2) / 60
}

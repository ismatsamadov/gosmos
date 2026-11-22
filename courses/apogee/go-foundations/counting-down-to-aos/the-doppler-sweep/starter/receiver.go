package main

import "math"

// signalStrength returns the beacon's received power in dBm with the receiver tuned offsetHz from the prediction.
func signalStrength(offsetHz int) float64 {
	return -112 - math.Abs(float64(offsetHz-3500))/500
}

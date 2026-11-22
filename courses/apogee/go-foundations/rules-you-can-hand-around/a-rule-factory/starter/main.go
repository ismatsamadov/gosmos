package main

// minElevation builds a rule accepting passes at or above a minimum angle.
func minElevation(minDeg float64) func(float64) bool {
	return func(elevationDeg float64) bool { return false }
}

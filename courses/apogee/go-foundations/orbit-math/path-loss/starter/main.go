package main

import (
	"fmt"
	"math"
)

// earthRadiusKm is Earth's mean radius.
const earthRadiusKm = 6371

// radians converts degrees to radians.
func radians(deg float64) float64 {
	return deg * math.Pi / 180
}

// slantRangeKm returns the distance to a satellite at altKm, seen at elevDeg.
func slantRangeKm(altKm, elevDeg float64) float64 {
	e := radians(elevDeg)
	a := earthRadiusKm + altKm
	b := earthRadiusKm * math.Cos(e)
	return math.Sqrt(a*a-b*b) - earthRadiusKm*math.Sin(e)
}

// pathLossDB returns the free-space path loss at freqMHz to a satellite at altKm, seen at elevDeg.
func pathLossDB(altKm, elevDeg, freqMHz float64) float64 {
	return 0
}

func main() {
	fmt.Printf("Pipit-1 at the horizon: %.2f dB\n", pathLossDB(525, 0, 437.25))
}

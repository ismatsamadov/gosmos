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

func main() {
	fmt.Printf("Pipit-1 at 30 deg: %.2f km\n", slantRangeKm(525, 30))
}

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
	return 0
}

func main() {
	fmt.Printf("Pipit-1 at 30 deg: %.2f km\n", slantRangeKm(525, 30))
}

package main

import "fmt"

// earthRadiusKm is Earth's mean radius.
const earthRadiusKm = 6371

// muKm3S2 is Earth's gravitational parameter, μ, in km³/s².
const muKm3S2 = 398600.4418

// periodMinutes returns the period of a circular orbit altKm above the Earth, in minutes.
func periodMinutes(altKm float64) float64 {
	return 0
}

func main() {
	fmt.Println("Pipit-1 period, min:", periodMinutes(525))
}

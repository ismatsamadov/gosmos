package main

import "fmt"

func main() {
	var name string = "Pipit-1"
	var altitudeKm int = 525
	var beaconMHz float64 = 437.25
	var deployed bool = true

	fmt.Println("Name:", name)
	fmt.Println("Altitude km:", altitudeKm)
	fmt.Println("Beacon MHz:", beaconMHz)
	fmt.Println("Deployed:", deployed)
}

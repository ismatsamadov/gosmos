package main

import "fmt"

const (
	fullMV   = 8200 // where the FULL band starts
	chargeMV = 150  // what Pipit-3 gains in one orbit
)

func main() {
	battMV := 6700
	orbits := 0
	for battMV < fullMV {
		battMV += chargeMV
		orbits++
	}
	fmt.Println("Orbits to full:", orbits)
	fmt.Println("Battery:", battMV, "mV")
}

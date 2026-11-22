package main

import "fmt"

// Station is one of Apogee's ground stations.
type Station struct {
	Code   string
	Name   string
	Stowed bool
}

// stow parks the station's dishes.
func stow(s *Station) {
	s.Stowed = true
}

func main() {
	yst := Station{Code: "YST", Name: "Ystmark"}
	stow(&yst)
	fmt.Println("YST stowed:", yst.Stowed)
}

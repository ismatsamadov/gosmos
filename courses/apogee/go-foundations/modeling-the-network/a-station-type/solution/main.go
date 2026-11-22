package main

import "fmt"

// Station is one of Apogee's ground stations.
type Station struct {
	Code    string
	Name    string
	LatDeg  float64 // north is positive
	MaskDeg float64
}

func main() {
	pel := Station{Code: "PEL", Name: "Port Ellery", LatDeg: 46.9, MaskDeg: 5}
	fmt.Printf("%+v\n", pel)
}

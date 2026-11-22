package main

import "fmt"

// Station is one of Apogee's ground stations.
type Station struct {
	Code    string
	Name    string
	LatDeg  float64 // north is positive
	MaskDeg float64
}

// stationLine formats a station as one line for the Board.
func stationLine(s Station) string {
	return ""
}

func main() {
	sqt := Station{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, MaskDeg: 3}
	fmt.Println(stationLine(sqt))
}

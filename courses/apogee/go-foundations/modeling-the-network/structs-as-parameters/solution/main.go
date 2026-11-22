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
	hemi := "N"
	lat := s.LatDeg
	if lat < 0 {
		hemi = "S"
		lat = -lat
	}
	return fmt.Sprintf("%s %s (%.1f %s, mask %.1f)", s.Code, s.Name, lat, hemi, s.MaskDeg)
}

func main() {
	sqt := Station{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, MaskDeg: 3}
	fmt.Println(stationLine(sqt))
}

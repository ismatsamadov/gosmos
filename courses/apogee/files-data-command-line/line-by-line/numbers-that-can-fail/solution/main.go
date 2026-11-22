package main

import "strconv"

// parseElevation reads an elevation in degrees, refusing text that is not a number.
func parseElevation(text string) (float64, error) {
	elevDeg, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}
	return elevDeg, nil
}

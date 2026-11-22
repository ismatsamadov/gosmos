package main

import (
	"fmt"
	"strconv"
	"strings"
)

// parseColumns splits one reading into its four columns: time, dish, azimuth and elevation.
func parseColumns(line string) ([]string, error) {
	cols := strings.Fields(line)
	if len(cols) != 4 {
		return nil, fmt.Errorf("%d columns instead of 4", len(cols))
	}
	return cols, nil
}

// parseElevation reads an elevation in degrees, refusing text that is not a number.
func parseElevation(text string) (float64, error) {
	elevDeg, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}
	return elevDeg, nil
}

// elevationOf returns the elevation of one reading, or why the line is not a valid reading.
func elevationOf(line string) (float64, error) {
	cols, err := parseColumns(line)
	if err != nil {
		return 0, err
	}
	return parseElevation(cols[3])
}

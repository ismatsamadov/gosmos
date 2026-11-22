package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Summary is the controller log in a few numbers.
type Summary struct {
	Readings   int            // readings in the whole log
	ByDish     map[string]int // readings per dish: "LG" for La Grande, "CH" for Chiquita
	MaxElevDeg float64        // the highest elevation of any reading; 0 when there are none
}

// isReading reports whether a line is a reading rather than a note or a blank line.
func isReading(line string) bool {
	view := strings.TrimLeft(line, " \t")
	return view != "" && !strings.HasPrefix(view, "#")
}

// parseReading returns a reading's dish and elevation, or why the line is not a valid reading.
func parseReading(line string) (dish string, elevDeg float64, err error) {
	cols := strings.Fields(line)
	if len(cols) != 4 {
		return "", 0, fmt.Errorf("%d columns instead of 4", len(cols))
	}
	elevDeg, err = strconv.ParseFloat(cols[3], 64)
	if err != nil {
		return "", 0, err
	}
	return cols[1], elevDeg, nil
}

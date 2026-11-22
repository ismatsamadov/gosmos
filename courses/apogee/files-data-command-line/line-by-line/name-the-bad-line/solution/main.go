package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// readElevations returns the elevation of every reading in the log, naming the line of the first bad one.
func readElevations(r io.Reader) ([]float64, error) {
	sc := bufio.NewScanner(r)
	var elevs []float64
	n := 0
	for sc.Scan() {
		n++
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue // a blank line between shifts, or a note
		}
		elevDeg, err := elevationOf(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", n, err)
		}
		elevs = append(elevs, elevDeg)
	}
	return elevs, sc.Err()
}

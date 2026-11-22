package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// readElevations returns the elevation of every reading in the log, naming the line of the first bad one.
func readElevations(r io.Reader) ([]float64, error) {
	_ = fmt.Errorf
	sc := bufio.NewScanner(r)
	var elevs []float64
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue // a blank line between shifts, or a note
		}
		elevDeg, err := elevationOf(line)
		if err != nil {
			return nil, err
		}
		elevs = append(elevs, elevDeg)
	}
	return elevs, sc.Err()
}

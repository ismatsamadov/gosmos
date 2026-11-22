package main

import (
	"bufio"
	"fmt"
	"io"
)

// summarize reads the log one line at a time into a Summary, naming the line of any bad reading.
func summarize(r io.Reader) (Summary, error) {
	s := Summary{ByDish: map[string]int{}}
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {
		n++
		if !isReading(sc.Text()) {
			continue
		}
		dish, elevDeg, err := parseReading(sc.Text())
		if err != nil {
			return Summary{}, fmt.Errorf("line %d: %w", n, err)
		}
		if s.Readings == 0 || elevDeg > s.MaxElevDeg {
			s.MaxElevDeg = elevDeg
		}
		s.Readings++
		s.ByDish[dish]++
	}
	return s, sc.Err()
}

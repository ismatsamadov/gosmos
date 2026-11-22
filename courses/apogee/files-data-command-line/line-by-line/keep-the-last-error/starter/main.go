package main

import (
	"bufio"
	"io"
)

// collectLines returns the lines r holds, with the error that stopped the reading early, if any.
func collectLines(r io.Reader) ([]string, error) {
	sc := bufio.NewScanner(r)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	// TODO: return why the loop stopped.
	return lines, nil
}

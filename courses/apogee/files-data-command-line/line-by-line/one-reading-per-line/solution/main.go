package main

import (
	"bufio"
	"io"
)

// countLines counts the lines r holds, blank ones and a last line without a newline included.
func countLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}

package main

import (
	"bufio"
	"io"
)

// countLines counts the lines r holds, blank ones and a last line without a newline included.
func countLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	n := 0
	// TODO: count each line sc.Scan finds.
	return n, sc.Err()
}

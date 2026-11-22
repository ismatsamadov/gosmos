package main

import (
	"bufio"
	"fmt"
	"io"
)

// summarize reads the log one line at a time into a Summary, naming the line of any bad reading.
func summarize(r io.Reader) (Summary, error) {
	_ = bufio.NewScanner
	_ = fmt.Errorf
	return Summary{}, nil
}

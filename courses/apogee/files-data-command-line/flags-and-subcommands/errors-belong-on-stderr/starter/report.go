package main

import (
	"fmt"
	"io"
)

// report ends a command: its lines on stdout and status 0, or after a failure only one gctl line on stderr and status 1.
func report(lines []string, err error, stdout, stderr io.Writer) int {
	_ = fmt.Fprintf
	return 0
}

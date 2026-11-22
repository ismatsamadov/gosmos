package main

import (
	"fmt"
	"io"
)

// report ends a command: its lines on stdout and status 0, or after a failure only one gctl line on stderr and status 1.
func report(lines []string, err error, stdout, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintf(stderr, "gctl: %v\n", err)
		return 1
	}
	for _, line := range lines {
		fmt.Fprintln(stdout, line)
	}
	return 0
}

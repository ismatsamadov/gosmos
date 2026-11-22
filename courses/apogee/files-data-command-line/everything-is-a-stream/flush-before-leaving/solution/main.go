package main

import (
	"bufio"
	"fmt"
	"io"
)

// bufferedReport writes every line and flushes it before returning.
func bufferedReport(w io.Writer, lines []string) error {
	buffer := bufio.NewWriter(w)
	for _, line := range lines {
		if _, err := fmt.Fprintln(buffer, line); err != nil {
			return err
		}
	}
	return buffer.Flush()
}

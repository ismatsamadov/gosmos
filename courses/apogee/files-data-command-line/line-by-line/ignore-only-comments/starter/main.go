package main

import (
	"bufio"
	"io"
	"strings"
)

// dataLines returns the log's readings exactly as written, leaving out notes and blank lines.
func dataLines(r io.Reader) ([]string, error) {
	sc := bufio.NewScanner(r)
	var lines []string
	for sc.Scan() {
		line := sc.Text()
		// TODO: classify by the first character that is not a space or tab.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, sc.Err()
}

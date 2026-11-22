package main

import "fmt"

// timeline lays spans out for the review: each root by its start, then its children in the order
// they started, every line with its start time and how long the span took.
func timeline(spans []Span) []string {
	var lines []string
	for _, s := range spans { // TODO: roots by start, each followed by its children, in the review's layout
		lines = append(lines, fmt.Sprintf("%s %s", s.Start.Format("15:04:05"), s.Name))
	}
	return lines
}

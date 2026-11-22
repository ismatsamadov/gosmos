package main

import (
	"fmt"
	"slices"
)

// timeline lays spans out for the review: each root by its start, then its children in the order
// they started, every line with its start time and how long the span took.
func timeline(spans []Span) []string {
	byStart := slices.Clone(spans)
	slices.SortStableFunc(byStart, func(a, b Span) int { return a.Start.Compare(b.Start) })
	var lines []string
	for _, root := range byStart {
		if root.ParentID != "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s %-24s %8s", root.Start.Format("15:04:05"), "pass "+root.PassID, root.End.Sub(root.Start)))
		for _, s := range byStart {
			if s.ParentID == root.SpanID {
				lines = append(lines, fmt.Sprintf("%s   %-22s %8s", s.Start.Format("15:04:05"), s.Name, s.End.Sub(s.Start)))
			}
		}
	}
	return lines
}

package main

import (
	"fmt"
	"maps"
	"slices"
)

// report returns one "id: count" line per satellite, sorted by satellite ID.
func report(counts map[string]int) []string {
	var lines []string
	for _, id := range slices.Sorted(maps.Keys(counts)) {
		lines = append(lines, fmt.Sprintf("%s: %d", id, counts[id]))
	}
	return lines
}

func main() {
	counts := map[string]int{"thistle": 2, "pipit-1": 3, "brine-3": 1}
	for _, line := range report(counts) {
		fmt.Println(line)
	}
}

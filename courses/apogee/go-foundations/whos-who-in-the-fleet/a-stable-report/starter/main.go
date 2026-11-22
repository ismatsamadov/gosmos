package main

import "fmt"

// report returns one "id: count" line per satellite, sorted by satellite ID.
func report(counts map[string]int) []string {
	return nil
}

func main() {
	counts := map[string]int{"thistle": 2, "pipit-1": 3, "brine-3": 1}
	for _, line := range report(counts) {
		fmt.Println(line)
	}
}

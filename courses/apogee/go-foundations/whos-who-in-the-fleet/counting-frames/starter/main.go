package main

import "fmt"

// countBySat counts the frames each satellite sent, by satellite ID.
func countBySat(ids []string) map[string]int {
	counts := map[string]int{}
	return counts
}

func main() {
	frames := []string{"pipit-1", "thistle", "pipit-1", "brine-3", "thistle", "pipit-1"}
	fmt.Println(countBySat(frames))
}

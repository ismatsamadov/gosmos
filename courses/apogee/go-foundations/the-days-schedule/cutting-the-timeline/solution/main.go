package main

import "fmt"

// nextN returns the first n pass IDs, or all of them when there are fewer.
func nextN(ids []string, n int) []string {
	n = min(n, len(ids))
	return ids[:n]
}

// lastN returns the last n readings, or all of them when there are fewer.
func lastN(readings []float64, n int) []float64 {
	n = min(n, len(readings))
	return readings[len(readings)-n:]
}

func main() {
	ids := []string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218", "cirrine-1/YST/07455", "pipit-3/YST/00214"}
	fmt.Println(nextN(ids, 3))
	fmt.Println(lastN([]float64{7.9, 7.8, 7.7}, 2))
}

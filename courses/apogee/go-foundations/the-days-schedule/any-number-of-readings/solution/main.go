package main

import "fmt"

// average returns the mean of the readings, or 0 when there are none.
func average(readings ...float64) float64 {
	if len(readings) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range readings {
		sum += r
	}
	return sum / float64(len(readings))
}

func main() {
	fmt.Println(average(-4.5, -3.9, -2.4))
}

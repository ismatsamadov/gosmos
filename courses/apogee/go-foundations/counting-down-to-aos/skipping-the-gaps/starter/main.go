package main

import "fmt"

// averageBattery returns the mean reading from minute from up to, not including, minute to, skipping lost frames.
func averageBattery(from, to int) float64 {
	return 0
}

func main() {
	fmt.Printf("Average battery: %.2f V\n", averageBattery(0, 6))
}

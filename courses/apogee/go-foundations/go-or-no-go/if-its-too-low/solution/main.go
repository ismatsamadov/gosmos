package main

import "fmt"

func main() {
	maxElevDeg := 42.0
	fmt.Printf("Pass 1: max %.1f deg\n", maxElevDeg)
	if maxElevDeg < 10 {
		fmt.Println("  low pass: expect a weak signal")
	}

	maxElevDeg = 7.5
	fmt.Printf("Pass 2: max %.1f deg\n", maxElevDeg)
	if maxElevDeg < 10 {
		fmt.Println("  low pass: expect a weak signal")
	}

	maxElevDeg = 15.2
	fmt.Printf("Pass 3: max %.1f deg\n", maxElevDeg)
	if maxElevDeg < 10 {
		fmt.Println("  low pass: expect a weak signal")
	}
}

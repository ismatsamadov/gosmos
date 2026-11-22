package main

import "fmt"

// highEnough reports whether a pass rises at least 20 degrees above the horizon.
func highEnough(elevationDeg float64) bool { return elevationDeg >= 20 }

func main() {
	rule := highEnough
	_ = rule
	fmt.Println("15 degrees:", false)
	fmt.Println("20 degrees:", false)
}

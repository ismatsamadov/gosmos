package main

import "fmt"

// highEnough reports whether a pass rises at least 20 degrees above the horizon.
func highEnough(elevationDeg float64) bool { return elevationDeg >= 20 }

func main() {
	rule := highEnough
	fmt.Println("15 degrees:", rule(15))
	fmt.Println("20 degrees:", rule(20))
}

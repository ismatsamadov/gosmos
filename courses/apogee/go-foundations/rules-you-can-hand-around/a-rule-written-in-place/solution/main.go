package main

import "fmt"

func main() {
	rule := func(elevationDeg float64) bool { return elevationDeg >= 20 }
	fmt.Println("19 degrees:", rule(19))
	fmt.Println("20 degrees:", rule(20))
	fmt.Println("45 degrees:", rule(45))
}

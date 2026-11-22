package main

import "fmt"

// Orbit periods in minutes, and the minutes in a day.
const (
	pipit1PeriodMin   = 94.985
	thistlePeriodMin  = 92.722
	cirrine1PeriodMin = 96.539
	minutesPerDay     = 1440
)

// TODO: declare orbitsPerDay here.

func main() {
	fmt.Printf("Pipit-1 %.2f\n", 15.0)
	fmt.Printf("Thistle %.2f\n", 15.4)
	fmt.Printf("Cirrine-1 %.2f\n", 14.9)
}

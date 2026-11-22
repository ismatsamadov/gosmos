package main

import "fmt"

// Orbit periods in minutes, and the minutes in a day.
const (
	pipit1PeriodMin   = 94.985
	thistlePeriodMin  = 92.722
	cirrine1PeriodMin = 96.539
	minutesPerDay     = 1440
)

// orbitsPerDay returns how many orbits of periodMin minutes fit in one day.
func orbitsPerDay(periodMin float64) float64 {
	return minutesPerDay / periodMin
}

func main() {
	fmt.Printf("Pipit-1 %.2f\n", orbitsPerDay(pipit1PeriodMin))
	fmt.Printf("Thistle %.2f\n", orbitsPerDay(thistlePeriodMin))
	fmt.Printf("Cirrine-1 %.2f\n", orbitsPerDay(cirrine1PeriodMin))
}

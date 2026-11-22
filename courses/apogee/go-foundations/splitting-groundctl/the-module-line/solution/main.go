package main

import (
	"fmt"

	"apogee/groundctl/orbit"
)

func main() {
	fmt.Printf("Pipit-1 period: %.2f min\n", orbit.PeriodMinutes(525))
}

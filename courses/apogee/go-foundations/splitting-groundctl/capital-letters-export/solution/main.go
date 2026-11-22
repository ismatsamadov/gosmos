package main

import (
	"fmt"

	"apogee/groundctl/orbit"
)

func main() {
	fmt.Printf("Thistle period: %.2f min\n", orbit.PeriodMinutes(415))
}

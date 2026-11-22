package main

import "fmt"

// printDishes prints Ystmark's dish and the test van's, one line each.
func printDishes(yst Station, van Van) {
	fmt.Printf("%s dish %.1f m, X-band %v\n", yst.Code, yst.DiameterM, yst.XBand)
	fmt.Printf("van dish %.1f m, X-band %v\n", van.DiameterM, van.XBand)
}

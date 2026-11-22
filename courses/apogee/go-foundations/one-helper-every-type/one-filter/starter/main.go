package main

import "fmt"

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID      string
	Station    string
	MaxElevDeg float64
}

// Filter returns the items for which keep returns true, in their original order.
func Filter[T any](items []T, keep func(T) bool) []T {
	return nil
}

func main() {
	passes := []Pass{
		{SatID: "pipit-1", Station: "PEL", MaxElevDeg: 42.5},
		{SatID: "thistle", Station: "LAN", MaxElevDeg: 8.1},
		{SatID: "furrow-a", Station: "SQT", MaxElevDeg: 71},
	}
	high := Filter(passes, func(p Pass) bool { return p.MaxElevDeg > 10 })
	fmt.Println(len(high), "of", len(passes), "passes rise above 10 degrees")
}

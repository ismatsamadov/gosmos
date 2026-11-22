package main

import (
	"errors"
	"fmt"
)

// Target is the share of booked passes a station promises to fly, in basis points: 9950 is 99.50%.
type Target int

// The targets agreed on Monday 8 January 2029.
const (
	Standard Target = 9950 // every station's
	Atoll    Target = 9900 // Lantern Atoll's: its palms and downpours alone cost about one pass in a hundred
)

// availability is the share of booked passes that were flown: flown/booked, from 0 to 1.
func availability(flown, booked int) (float64, error) {
	if booked <= 0 {
		return 0, errors.New("no booked passes")
	}
	if flown > booked {
		return 0, errors.New("more flown than booked")
	}
	return float64(flown) / float64(booked), nil
}

// budget is how many of booked passes may be lost while t is still met, in whole passes.
func budget(t Target, booked int) int {
	return booked * (10000 - int(t)) / 10000
}

// row is one station's line of the quarter's report: code, booked, lost, availability in percent,
// budget and what is left of it.
func row(code string, booked, lost int, t Target) string {
	avail, err := availability(booked-lost, booked)
	if err != nil {
		return fmt.Sprintf("%-7s %v", code, err)
	}
	b := budget(t, booked)
	return fmt.Sprintf("%-7s %7d %5d %6.2f %7d %5d", code, booked, lost, avail*100, b, b-lost)
}

package main

// Target is the share of booked passes a station promises to fly, in basis points: 9950 is 99.50%.
type Target int

// The targets agreed on Monday 8 January 2029.
const (
	Standard Target = 9950 // every station's
	Atoll    Target = 9900 // Lantern Atoll's: its palms and downpours alone cost about one pass in a hundred
)

// budget is how many of booked passes may be lost while t is still met, in whole passes.
func budget(t Target, booked int) int {
	return booked * (10000 - int(t)) / 10000
}

package main

// Target is the share of booked passes a station promises to fly, in basis points: 9950 is 99.50%.
type Target int

// The targets agreed on Monday 8 January 2029.
const (
	Standard Target = 9950 // every station's
	Atoll    Target = 9900 // Lantern Atoll's: its palms and downpours alone cost about one pass in a hundred
)

// burnRate is how fast lost passes spend t's budget: at 1 it lasts exactly the window, at 50 a fiftieth of it.
func burnRate(t Target, booked, lost int) float64 {
	return float64(lost) / float64(booked)
}

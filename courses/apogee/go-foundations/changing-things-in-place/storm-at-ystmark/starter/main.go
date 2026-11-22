package main

import "fmt"

// stowKmh is the wind speed from which Ystmark's dishes stow.
const stowKmh = 72

// Dish is one of Ystmark's two antennas.
type Dish struct {
	Number int
	Stowed bool
	Stows  int // how many times the dish has been stowed
}

// Stow parks the dish, counting each new stow.
func (d *Dish) Stow() {
	if !d.Stowed {
		d.Stowed = true
		d.Stows++
	}
}

// Unstow sends the dish back to tracking.
func (d *Dish) Unstow() {
	d.Stowed = false
}

// react stows every dish when the wind reaches stowKmh, and unstows them below it.
func react(dishes []Dish, windKmh int) {
	for _, d := range dishes {
		if windKmh >= stowKmh {
			d.Stow()
		} else {
			d.Unstow()
		}
	}
}

// state reports "stowed" when every dish is stowed, and "tracking" otherwise.
func state(dishes []Dish) string {
	for _, d := range dishes {
		if !d.Stowed {
			return "tracking"
		}
	}
	return "stowed"
}

func main() {
	dishes := []Dish{{Number: 1}, {Number: 2}}
	winds := []int{64, 75, 58, 81} // km/h, one reading an hour from 02:00 UTC
	for i, wind := range winds {
		react(dishes, wind)
		fmt.Printf("%02d:00 wind %d km/h: %s\n", 2+i, wind, state(dishes))
	}
	for _, d := range dishes {
		fmt.Printf("dish %d stowed %d times\n", d.Number, d.Stows)
	}
}

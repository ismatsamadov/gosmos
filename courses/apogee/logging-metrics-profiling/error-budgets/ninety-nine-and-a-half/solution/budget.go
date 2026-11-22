package main

import "errors"

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

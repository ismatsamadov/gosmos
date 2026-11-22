package main

import "errors"

// availability is the share of booked passes that were flown: flown/booked, from 0 to 1.
func availability(flown, booked int) (float64, error) {
	_ = errors.New // TODO: refuse counts that can't be, then divide
	return 0, nil
}

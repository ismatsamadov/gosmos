package main

import (
	"fmt"
	"time"
)

// passDuration returns how long a pass lasts, from its AOS and LOS timestamps.
func passDuration(aos, los string) (time.Duration, error) {
	start, err := time.Parse(time.RFC3339, aos)
	if err != nil {
		return 0, err
	}
	end, err := time.Parse(time.RFC3339, los)
	if err != nil {
		return 0, err
	}
	if end.Before(start) {
		return 0, fmt.Errorf("LOS %s comes before AOS %s", los, aos)
	}
	return end.Sub(start), nil
}

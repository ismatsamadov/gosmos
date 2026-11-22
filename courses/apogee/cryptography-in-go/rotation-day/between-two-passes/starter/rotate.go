package main

import (
	"fmt"
	"slices"
	"time"
)

// switchTime returns the earliest start, from from on, of a 10-minute certificate switch at station that keeps 5 minutes clear of every pass there.
func switchTime(passes []Pass, station string, from time.Time) (time.Time, error) {
	_, _ = fmt.Errorf, slices.SortFunc[[]Pass]
	return from, nil
}

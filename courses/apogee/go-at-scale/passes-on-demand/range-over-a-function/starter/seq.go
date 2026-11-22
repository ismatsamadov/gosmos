package main

import (
	"iter"
	"time"
)

// Days yields n midnights, a day apart, starting at from.
func Days(from time.Time, n int) iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
	}
}

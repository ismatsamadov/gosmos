package main

import "time"

// exhaustIn is how long remaining passes of budget last at lostLast7d lost every 7 days; false when
// nothing was lost in the week, so there is no rate.
func exhaustIn(remaining, lostLast7d int) (time.Duration, bool) {
	if remaining <= 0 {
		return 0, true
	}
	if lostLast7d == 0 {
		return 0, false
	}
	return time.Duration(remaining) * 7 * 24 * time.Hour / time.Duration(lostLast7d), true
}

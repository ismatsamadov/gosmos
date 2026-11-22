package main

import "time"

// exhaustIn is how long remaining passes of budget last at lostLast7d lost every 7 days; false when
// nothing was lost in the week, so there is no rate.
func exhaustIn(remaining, lostLast7d int) (time.Duration, bool) {
	return 0, false
}

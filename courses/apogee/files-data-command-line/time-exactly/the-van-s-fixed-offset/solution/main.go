package main

import "time"

// vanTime returns t shown at a fixed offset of secondsEast seconds east of UTC.
func vanTime(t time.Time, secondsEast int) time.Time {
	return t.In(time.FixedZone("", secondsEast))
}

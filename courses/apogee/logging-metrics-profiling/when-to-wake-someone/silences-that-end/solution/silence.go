package main

import "time"

// Silence holds back a station's pages for a while. It always ends.
type Silence struct {
	Station string
	From    time.Time // included
	Until   time.Time // excluded
}

// silenced reports whether a silence in s covers station at the time at: From included,
// Until not.
func silenced(s []Silence, station string, at time.Time) bool {
	for _, x := range s {
		if x.Station == station && !at.Before(x.From) && at.Before(x.Until) {
			return true
		}
	}
	return false
}

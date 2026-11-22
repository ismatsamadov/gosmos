package main

import "time"

// Result is how one pass ended, as a station's uplinkd reports it to the pager.
type Result struct {
	PassID  string    // "cirrine-1/LAN/16402"
	Station string    // "LAN"
	LOS     time.Time // a result counts from its LOS
	Outcome string    // "flown", "failed" or "skipped"
	Reason  string    // "" when flown; "deadline", "no_ack", "no_worker", "radio"; "stowed" when skipped
}

// failedWithin counts station's failed passes whose LOS is in (now-d, now].
func failedWithin(results []Result, station string, now time.Time, d time.Duration) int {
	n := 0
	for _, r := range results {
		if r.Station == station && r.Outcome == "failed" && r.LOS.After(now.Add(-d)) && !r.LOS.After(now) {
			n++
		}
	}
	return n
}

package main

import "time"

// window counts station's results with LOS in the 30 days to now, (now − 30 days, now]: the passes
// booked, and those lost, failed or skipped.
func window(results []Result, station string, now time.Time) (booked, lost int) {
	for _, r := range results {
		booked++
		if r.Outcome != "flown" {
			lost++
		}
	}
	return booked, lost
}

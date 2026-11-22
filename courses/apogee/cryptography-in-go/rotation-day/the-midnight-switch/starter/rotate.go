package main

import "time"

// accepts reports whether a command sealed with the key version id is accepted at now.
func accepts(ring []KeyVersion, id string, now time.Time) bool {
	for i, v := range ring {
		if v.ID != id {
			continue
		}
		// the draft's hard switch: a version counts from its Activate until the next one's
		taken := i+1 < len(ring) && !now.Before(ring[i+1].Activate)
		return !now.Before(v.Activate) && !taken
	}
	return false
}

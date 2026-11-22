package main

import (
	"slices"
	"time"
)

// planStationLean plans s's candidates as planStationFat does, with
// less garbage.
func planStationLean(s Site, cands []PassView, turnaround time.Duration) []Contact {
	work := getScratch()
	defer putScratch(work)
	seen := make(map[passKey]bool, len(cands))
	for _, p := range cands {
		key := passKey{p.SatID, p.Station, p.AOS.Unix()}
		if seen[key] {
			continue
		}
		seen[key] = true
		*work = append(*work, p)
	}
	slices.SortFunc(*work, byAOS)
	contacts, _ := assign(s, slices.Values(*work), nil, turnaround)
	return contacts
}

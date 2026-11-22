package main

import (
	"fmt"
	"slices"
	"time"
)

// planStationLean plans s's candidates as planStationFat does, with
// less garbage.
func planStationLean(s Site, cands []PassView, turnaround time.Duration) []Contact {
	var work []PassView
	seen := map[string]bool{}
	for _, p := range cands {
		key := fmt.Sprintf("%s %s %d", p.SatID, p.Station, p.AOS.Unix())
		if seen[key] {
			continue
		}
		seen[key] = true
		work = append(work, p)
	}
	slices.SortFunc(work, byAOS)
	contacts, _ := assign(s, slices.Values(work), nil, turnaround)
	return contacts
}

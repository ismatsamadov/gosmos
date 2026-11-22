package main

import (
	"fmt"
	"slices"
	"time"
)

// planStationFat plans s's candidates by the rule, as passplan 2's
// draft does: it copies them, skipping any pass seen before, sorts
// the copy by AOS and pass ID, and gives its passes dishes.
func planStationFat(s Site, cands []PassView, turnaround time.Duration) []Contact {
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

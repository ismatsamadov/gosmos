package main

import (
	"cmp"
	"slices"
	"strings"
	"time"
)

// firstWaveDay returns Monday 2 April, the first wave's first day, as
// passplan 2 plans it from scratch, every dish free at midnight: 260
// contacts at the ten stations live that day, by AOS then pass ID. No
// two candidates at one station rise in the same second that day.
func firstWaveDay() *Schedule {
	n := newNetwork()
	from := time.Date(2029, 4, 2, 0, 0, 0, 0, time.UTC)
	sch := &Schedule{Window: Window{From: from, To: from.AddDate(0, 0, 1)}}
	for _, s := range n.Sites {
		var cands []PassView
		for i, e := range n.Elements {
			if !canServe(s, n.Satellites[i]) {
				continue
			}
			for _, p := range predict(e, s, sch.From, sch.To) {
				if p.MaxElevDeg >= 10 {
					cands = append(cands, p)
				}
			}
		}
		slices.SortFunc(cands, byAOS)
		contacts, dropped := assign(s, slices.Values(cands), nil, 2*time.Minute)
		sch.Contacts = append(sch.Contacts, contacts...)
		sch.Dropped += dropped
	}
	slices.SortFunc(sch.Contacts, func(a, b Contact) int {
		return cmp.Or(byAOS(a.PassView, b.PassView), strings.Compare(a.ID, b.ID))
	})
	return sch
}

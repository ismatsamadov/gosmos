package main

import (
	"slices"
	"time"
)

// takeOverDraft returns Thursday 28 June as passplan 2 plans it from scratch, every dish free
// at midnight: 2,020 contacts at twenty stations, by AOS then pass ID.
func takeOverDraft() []Contact {
	n := newNetwork()
	from := time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)
	var all []Contact
	for _, s := range n.Sites {
		var cands []PassView
		for i, e := range n.Elements {
			if !canServe(s, n.Satellites[i]) {
				continue
			}
			for _, p := range predict(e, s, from, from.Add(24*time.Hour)) {
				if p.MaxElevDeg >= 10 {
					cands = append(cands, p)
				}
			}
		}
		slices.SortFunc(cands, byAOS)
		contacts, _ := assign(s, slices.Values(cands), nil, 2*time.Minute)
		all = append(all, contacts...)
	}
	slices.SortFunc(all, func(a, b Contact) int { return byAOS(a.PassView, b.PassView) })
	return all
}

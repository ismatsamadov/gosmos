package main

import (
	"cmp"
	"slices"
	"strings"
	"time"
)

// thursday returns Thursday 5 April as the two planners planned it from scratch, every dish
// free at midnight. Tirrin Rock's host parked both dishes from 07:30 to 08:24 to clear their
// ice, a stow only the old planner was told of. The old planner lists its contacts by AOS,
// then pass ID; passplan 2 lists the two that rise at 05:00:11, brine-3/SQT/16651 and
// furrow-a/VSK/23251, the other way round, since byAOS leaves their order open.
func thursday() (oldPlan, newPlan []Contact) {
	day := time.Date(2029, 4, 5, 0, 0, 0, 0, time.UTC)
	iceFrom, iceTo := day.Add(7*time.Hour+30*time.Minute), day.Add(8*time.Hour+24*time.Minute)
	n := newNetwork()
	for _, s := range n.Sites {
		var all, told []PassView
		for i, e := range n.Elements {
			if !canServe(s, n.Satellites[i]) {
				continue
			}
			for _, p := range predict(e, s, day, day.AddDate(0, 0, 1)) {
				if p.MaxElevDeg < 10 {
					continue
				}
				all = append(all, p)
				if s.Code != "TIR" || !p.AOS.Before(iceTo) || !iceFrom.Before(p.LOS) {
					told = append(told, p)
				}
			}
		}
		slices.SortFunc(all, byAOS)
		slices.SortFunc(told, byAOS)
		o, _ := assign(s, slices.Values(told), nil, 2*time.Minute)
		p, _ := assign(s, slices.Values(all), nil, 2*time.Minute)
		oldPlan, newPlan = append(oldPlan, o...), append(newPlan, p...)
	}
	byAOSThenID := func(a, b Contact) int {
		return cmp.Or(byAOS(a.PassView, b.PassView), strings.Compare(a.ID, b.ID))
	}
	slices.SortFunc(oldPlan, byAOSThenID)
	slices.SortFunc(newPlan, byAOSThenID)
	i := slices.IndexFunc(newPlan, func(c Contact) bool { return c.ID == "brine-3/SQT/16651" })
	newPlan[i], newPlan[i+1] = newPlan[i+1], newPlan[i]
	return oldPlan, newPlan
}

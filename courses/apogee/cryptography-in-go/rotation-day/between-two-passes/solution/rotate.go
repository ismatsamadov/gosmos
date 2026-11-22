package main

import (
	"fmt"
	"slices"
	"time"
)

// switchTime returns the earliest start, from from on, of a 10-minute certificate switch at station that keeps 5 minutes clear of every pass there.
func switchTime(passes []Pass, station string, from time.Time) (time.Time, error) {
	var mine []Pass
	for _, p := range passes {
		if p.Station == station {
			mine = append(mine, p)
		}
	}
	slices.SortFunc(mine, func(a, b Pass) int { return a.AOS.Compare(b.AOS) })
	start := from
	for _, p := range mine {
		if !start.Add(10 * time.Minute).After(p.AOS.Add(-5 * time.Minute)) {
			break // the switch ends before this pass's clear zone starts
		}
		if clear := p.LOS.Add(5 * time.Minute); start.Before(clear) {
			start = clear
		}
	}
	if end := from.Add(24 * time.Hour); start.Add(10 * time.Minute).After(end) {
		return time.Time{}, fmt.Errorf("no 10-minute gap for %s before %s", station, end.Format(time.RFC3339))
	}
	return start, nil
}

package main

import (
	"cmp"
	"fmt"
	"slices"
	"time"
)

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID   string
	Station string    // station code
	AOS     time.Time // acquisition of signal
	LOS     time.Time // loss of signal
}

// utc returns 21 April 2027 at h:m UTC.
func utc(h, m int) time.Time {
	return time.Date(2027, time.April, 21, h, m, 0, 0, time.UTC)
}

// sortByAOS puts passes in AOS order, and passes with the same AOS in station order.
func sortByAOS(passes []Pass) {
	slices.SortFunc(passes, func(a, b Pass) int {
		c := a.AOS.Compare(b.AOS)
		if c != 0 {
			return c
		}
		return cmp.Compare(a.Station, b.Station)
	})
}

func main() {
	passes := []Pass{
		{SatID: "cirrine-1", Station: "YST", AOS: utc(7, 30), LOS: utc(7, 41)},
		{SatID: "pipit-1", Station: "PEL", AOS: utc(4, 12), LOS: utc(4, 21)},
		{SatID: "furrow-a", Station: "SQT", AOS: utc(5, 3), LOS: utc(5, 14)},
	}
	sortByAOS(passes)
	for _, p := range passes {
		fmt.Println(p.SatID, p.Station)
	}
}

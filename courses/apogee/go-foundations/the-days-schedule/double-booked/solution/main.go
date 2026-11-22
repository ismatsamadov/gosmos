package main

import (
	"fmt"
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

// overlaps describes every pair of passes that need the same station at the same time.
func overlaps(passes []Pass) []string {
	var lines []string
	for i := range passes {
		for j := i + 1; j < len(passes); j++ {
			a, b := passes[i], passes[j]
			if a.Station == b.Station && a.AOS.Before(b.LOS) && b.AOS.Before(a.LOS) {
				lines = append(lines, fmt.Sprintf("%s@%s overlaps %s@%s", a.SatID, a.Station, b.SatID, b.Station))
			}
		}
	}
	return lines
}

func main() {
	tonight := []Pass{
		{SatID: "furrow-b", Station: "YST", AOS: utc(21, 40), LOS: utc(21, 52)},
		{SatID: "cirrine-1", Station: "SQT", AOS: utc(21, 44), LOS: utc(21, 55)},
		{SatID: "pipit-3", Station: "YST", AOS: utc(21, 49), LOS: utc(21, 59)},
	}
	fmt.Println(overlaps(tonight))
}

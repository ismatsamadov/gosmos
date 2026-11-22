package main

import (
	"fmt"
	"time"
)

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID   string
	Station string // station code
	AOS     time.Time
	LOS     time.Time
}

// Duration returns how long the pass lasts, from AOS to LOS.
func (p Pass) Duration() time.Duration {
	return p.LOS.Sub(p.AOS)
}

// at returns h:m:s UTC on 4 May 2027.
func at(h, m, s int) time.Time {
	return time.Date(2027, time.May, 4, h, m, s, 0, time.UTC)
}

func main() {
	passes := []Pass{
		{SatID: "pipit-1", Station: "PEL", AOS: at(4, 12, 0), LOS: at(4, 21, 12)},
		{SatID: "furrow-a", Station: "SQT", AOS: at(5, 3, 0), LOS: at(5, 14, 4)},
		{SatID: "thistle", Station: "LAN", AOS: at(6, 47, 0), LOS: at(6, 53, 30)},
	}
	for _, p := range passes {
		fmt.Printf("%s@%s %v\n", p.SatID, p.Station, p.Duration())
	}
}

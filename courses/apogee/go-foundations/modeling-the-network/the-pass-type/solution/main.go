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

// duration returns how long a pass lasts, from AOS to LOS.
func duration(p Pass) time.Duration {
	return p.LOS.Sub(p.AOS)
}

func main() {
	p := Pass{
		SatID:   "pipit-1",
		Station: "PEL",
		AOS:     time.Date(2027, 4, 15, 4, 12, 0, 0, time.UTC),
		LOS:     time.Date(2027, 4, 15, 4, 21, 12, 0, time.UTC),
	}
	fmt.Println(p.SatID, "over", p.Station, "for", duration(p))
}

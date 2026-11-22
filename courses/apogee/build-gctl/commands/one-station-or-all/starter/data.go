package main

import (
	"slices"
	"time"
)

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID      string
	Station    string    // station code
	AOS        time.Time // acquisition of signal
	LOS        time.Time // loss of signal
	MaxElevDeg float64
}

// Station is one of Apogee's ground stations.
type Station struct {
	Code       string // "YST": three capital letters
	Name       string // "Ystmark"
	Stowed     bool
	StowReason string
}

// tonight returns the night's passes in the schedule file's order, station by station. The
// night of Friday 29 October 2027 runs past midnight, so its passes are on 30 October, UTC.
func tonight() []Pass {
	at := func(h, m int) time.Time { return time.Date(2027, time.October, 30, h, m, 0, 0, time.UTC) }
	return []Pass{
		{"pipit-2", "YST", at(4, 31), at(4, 40), 47.2},
		{"cirrine-1", "YST", at(5, 47), at(5, 57), 71.9},
		{"pipit-1", "PEL", at(4, 12), at(4, 21), 61.4},
		{"brine-3", "LAN", at(5, 20), at(5, 26), 22.8},
		{"furrow-a", "CWR", at(4, 58), at(5, 6), 33.5},
	}
}

// stationTable returns Apogee's five ground stations, none of them stowed.
func stationTable() []Station {
	return []Station{
		{Code: "YST", Name: "Ystmark"},
		{Code: "PEL", Name: "Port Ellery"},
		{Code: "SQT", Name: "Salar Quieto"},
		{Code: "LAN", Name: "Lantern Atoll"},
		{Code: "CWR", Name: "Cape Wren"},
	}
}

// stationIndex returns the index of the station with the given code in stations, or -1.
func stationIndex(stations []Station, code string) int {
	return slices.IndexFunc(stations, func(st Station) bool { return st.Code == code })
}

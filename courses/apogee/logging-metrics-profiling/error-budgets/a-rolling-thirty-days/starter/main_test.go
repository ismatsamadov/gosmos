package main

import (
	"testing"
	"time"
)

// at parses a UTC time written in RFC 3339.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// TestWindowLanternAtoll replays Lantern Atoll's results from 8 December to the morning of
// 8 January: six natural losses and the five of 27-28 December are in its last 30 days.
func TestWindowLanternAtoll(t *testing.T) {
	booked, lost := window(lanResults(), "LAN", at("2029-01-08T00:00:00Z"))
	if booked != 387 || lost != 11 {
		t.Errorf(`window(results, "LAN", 2029-01-08T00:00:00Z) = %d, %d, want 387, 11`, booked, lost)
	}
}

func TestWindowEdges(t *testing.T) {
	month := 30 * 24 * time.Hour
	brine := Result{PassID: "brine-3/LAN/15145", Station: "LAN", LOS: at("2028-12-27T19:41:28Z"), Outcome: "failed", Reason: "deadline"}
	cape := []Result{
		{PassID: "cirrine-1/CWR/17147", Station: "CWR", LOS: at("2028-12-28T01:19:19Z"), Outcome: "flown"},
		{PassID: "furrow-a/CWR/21758", Station: "CWR", LOS: at("2028-12-28T02:30:56Z"), Outcome: "skipped", Reason: "stowed"},
	}
	tests := []struct {
		name         string
		results      []Result
		station      string
		now          time.Time
		booked, lost int
	}{
		{"Brine-3's failure, its LOS exactly now", []Result{brine}, "LAN", brine.LOS, 1, 1},
		{"30 days less a second after it", []Result{brine}, "LAN", brine.LOS.Add(month - time.Second), 1, 1},
		{"exactly 30 days after it", []Result{brine}, "LAN", brine.LOS.Add(month), 0, 0},
		{"a second before its LOS", []Result{brine}, "LAN", brine.LOS.Add(-time.Second), 0, 0},
		{"another station's window", []Result{brine}, "CWR", brine.LOS, 0, 0},
		{"Cape Wren stowed in wind on 28 December", cape, "CWR", at("2028-12-28T03:00:00Z"), 2, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booked, lost := window(tt.results, tt.station, tt.now)
			if booked != tt.booked || lost != tt.lost {
				t.Errorf("window(results, %q, %s) = %d, %d, want %d, %d", tt.station, tt.now.Format(time.RFC3339), booked, lost, tt.booked, tt.lost)
			}
		})
	}
}

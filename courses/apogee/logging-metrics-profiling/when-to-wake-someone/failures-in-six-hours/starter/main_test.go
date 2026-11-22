package main

import (
	"testing"
	"time"
)

// at returns hh:mm:ss UTC on the given day of December 2028.
func at(day, h, m, s int) time.Time {
	return time.Date(2028, 12, day, h, m, s, 0, time.UTC)
}

// night is how the passes of the night of 27-28 December ended: Lantern Atoll's six failures
// and two flown passes, and a Cape Wren pass skipped while its dish was stowed.
var night = []Result{
	{PassID: "brine-3/LAN/15145", Station: "LAN", LOS: at(27, 19, 41, 28), Outcome: "failed", Reason: "deadline"},
	{PassID: "pipit-1/LAN/14031", Station: "LAN", LOS: at(27, 20, 20, 49), Outcome: "failed", Reason: "deadline"},
	{PassID: "cirrine-1/LAN/17147", Station: "LAN", LOS: at(28, 0, 45, 16), Outcome: "failed", Reason: "deadline"},
	{PassID: "furrow-a/LAN/21758", Station: "LAN", LOS: at(28, 1, 57, 49), Outcome: "failed", Reason: "deadline"},
	{PassID: "cirrine-1/LAN/17148", Station: "LAN", LOS: at(28, 2, 19, 14), Outcome: "failed", Reason: "deadline"},
	{PassID: "furrow-a/CWR/21758", Station: "CWR", LOS: at(28, 2, 30, 56), Outcome: "skipped", Reason: "stowed"},
	{PassID: "furrow-b/LAN/21756", Station: "LAN", LOS: at(28, 2, 44, 57), Outcome: "failed", Reason: "deadline"},
	{PassID: "brine-3/LAN/15152", Station: "LAN", LOS: at(28, 7, 26, 4), Outcome: "flown"},
	{PassID: "pipit-1/LAN/14038", Station: "LAN", LOS: at(28, 8, 12, 34), Outcome: "flown"},
}

func TestFailedWithin(t *testing.T) {
	tests := []struct {
		name    string
		results []Result
		station string
		now     time.Time
		d       time.Duration
		want    int
	}{
		{"20:20:49, a LOS exactly at now counts", night, "LAN", at(27, 20, 20, 49), 6 * time.Hour, 2},
		{"01:41:28, the first failure exactly six hours back is out", night, "LAN", at(28, 1, 41, 28), 6 * time.Hour, 2},
		{"02:45, the worst of the night", night, "LAN", at(28, 2, 45, 0), 6 * time.Hour, 4},
		{"08:44, the last failure still inside", night, "LAN", at(28, 8, 44, 0), 6 * time.Hour, 1},
		{"08:45, flown passes don't count", night, "LAN", at(28, 8, 45, 0), 6 * time.Hour, 0},
		{"Ystmark, another station's failures don't count", night, "YST", at(28, 2, 45, 0), 6 * time.Hour, 0},
		{"Cape Wren, a skipped pass isn't a failure", night, "CWR", at(28, 3, 0, 0), 6 * time.Hour, 0},
		{"19:42, the last minute", night, "LAN", at(27, 19, 42, 0), time.Minute, 1},
		{"19:43, the minute after", night, "LAN", at(27, 19, 43, 0), time.Minute, 0},
		{"no results", nil, "LAN", at(28, 2, 45, 0), 6 * time.Hour, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			what := "the night"
			if tt.results == nil {
				what = "nil"
			}
			if got := failedWithin(tt.results, tt.station, tt.now, tt.d); got != tt.want {
				t.Errorf("failedWithin(%s, %q, %s, %v) = %d, want %d", what, tt.station, tt.now.Format("Jan 2 15:04:05"), tt.d, got, tt.want)
			}
		})
	}
}

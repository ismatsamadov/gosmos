package main

import (
	"slices"
	"testing"
	"time"
)

// at returns 21 April 2027 at h:m UTC.
func at(h, m int) time.Time {
	return time.Date(2027, time.April, 21, h, m, 0, 0, time.UTC)
}

// order lists passes as sat@station, in slice order.
func order(passes []Pass) []string {
	var out []string
	for _, p := range passes {
		out = append(out, p.SatID+"@"+p.Station)
	}
	return out
}

func TestSortByAOS(t *testing.T) {
	tests := []struct {
		name   string
		passes []Pass
		want   []string
	}{
		{"four passes", []Pass{
			{SatID: "cirrine-1", Station: "YST", AOS: at(7, 30), LOS: at(7, 41)},
			{SatID: "pipit-1", Station: "PEL", AOS: at(4, 12), LOS: at(4, 21)},
			{SatID: "thistle", Station: "LAN", AOS: at(6, 47), LOS: at(6, 53)},
			{SatID: "furrow-a", Station: "SQT", AOS: at(5, 3), LOS: at(5, 14)},
		}, []string{"pipit-1@PEL", "furrow-a@SQT", "thistle@LAN", "cirrine-1@YST"}},
		{"same AOS", []Pass{
			{SatID: "furrow-b", Station: "YST", AOS: at(21, 40), LOS: at(21, 52)},
			{SatID: "pipit-2", Station: "PEL", AOS: at(21, 15), LOS: at(21, 24)},
			{SatID: "brine-3", Station: "LAN", AOS: at(21, 40), LOS: at(21, 49)},
		}, []string{"pipit-2@PEL", "brine-3@LAN", "furrow-b@YST"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortByAOS(tt.passes)
			if got := order(tt.passes); !slices.Equal(got, tt.want) {
				t.Errorf("sortByAOS() gives %v, want %v", got, tt.want)
			}
		})
	}
}

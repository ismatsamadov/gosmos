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

func TestOverlaps(t *testing.T) {
	tests := []struct {
		name   string
		passes []Pass
		want   []string
	}{
		{"no overlap", []Pass{
			{SatID: "pipit-2", Station: "YST", AOS: at(20, 5), LOS: at(20, 15)},
			{SatID: "cirrine-1", Station: "YST", AOS: at(20, 30), LOS: at(20, 41)},
		}, nil},
		{"Furrow-B and Pipit-3 at Ystmark", []Pass{
			{SatID: "furrow-b", Station: "YST", AOS: at(21, 40), LOS: at(21, 52)},
			{SatID: "cirrine-1", Station: "SQT", AOS: at(21, 44), LOS: at(21, 55)},
			{SatID: "pipit-3", Station: "YST", AOS: at(21, 49), LOS: at(21, 59)},
		}, []string{"furrow-b@YST overlaps pipit-3@YST"}},
		{"touching passes", []Pass{
			{SatID: "furrow-b", Station: "YST", AOS: at(21, 40), LOS: at(21, 52)},
			{SatID: "pipit-3", Station: "YST", AOS: at(21, 52), LOS: at(22, 2)},
		}, nil},
		{"same times at different stations", []Pass{
			{SatID: "furrow-b", Station: "YST", AOS: at(21, 40), LOS: at(21, 52)},
			{SatID: "brine-3", Station: "LAN", AOS: at(21, 40), LOS: at(21, 52)},
		}, nil},
		{"three at once", []Pass{
			{SatID: "pipit-2", Station: "YST", AOS: at(22, 10), LOS: at(22, 20)},
			{SatID: "furrow-a", Station: "YST", AOS: at(22, 12), LOS: at(22, 23)},
			{SatID: "cirrine-1", Station: "YST", AOS: at(22, 15), LOS: at(22, 26)},
		}, []string{
			"pipit-2@YST overlaps furrow-a@YST",
			"pipit-2@YST overlaps cirrine-1@YST",
			"furrow-a@YST overlaps cirrine-1@YST",
		}},
		{"listed out of order", []Pass{
			{SatID: "pipit-3", Station: "YST", AOS: at(21, 49), LOS: at(21, 59)},
			{SatID: "pipit-2", Station: "YST", AOS: at(20, 5), LOS: at(20, 15)},
			{SatID: "furrow-b", Station: "YST", AOS: at(21, 40), LOS: at(21, 52)},
		}, []string{"pipit-3@YST overlaps furrow-b@YST"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := overlaps(tt.passes); !slices.Equal(got, tt.want) {
				t.Errorf("overlaps() = %q, want %q", got, tt.want)
			}
		})
	}
}

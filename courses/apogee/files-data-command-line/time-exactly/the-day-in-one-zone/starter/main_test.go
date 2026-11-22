package main

import (
	"slices"
	"testing"
	"time"
)

// at returns h:m:s UTC on the given day of September 2027.
func at(day, h, m, s int) time.Time {
	return time.Date(2027, time.September, day, h, m, s, 0, time.UTC)
}

// ids lists the bookings' IDs in order.
func ids(bookings []Booking) []string {
	var out []string
	for _, b := range bookings {
		out = append(out, b.ID)
	}
	return out
}

func TestDayBookings(t *testing.T) {
	week := []Booking{
		{ID: "bk-3104", SatID: "furrow-a", Station: "SQT", AOS: at(24, 6, 30, 0)},
		{ID: "bk-3101", SatID: "pipit-2", Station: "SQT", AOS: at(23, 15, 10, 0)},
		{ID: "bk-3099", SatID: "cirrine-1", Station: "SQT", AOS: at(23, 2, 40, 0)},
		{ID: "bk-3105", SatID: "brine-3", Station: "SQT", AOS: at(24, 0, 0, 0)},
		{ID: "bk-3102", SatID: "pipit-1", Station: "PEL", AOS: at(23, 7, 0, 0)},
		{ID: "bk-3100", SatID: "furrow-b", Station: "SQT", AOS: at(23, 6, 59, 59)},
	}
	van := time.FixedZone("", -7*60*60)
	for _, tc := range []struct {
		name     string
		day      string
		zoneName string
		zone     *time.Location
		want     []string
	}{
		{"the van's Thursday", "2027-09-23", "van", van, []string{"bk-3102", "bk-3101", "bk-3105", "bk-3104"}},
		{"the network's Thursday", "2027-09-23", "UTC", time.UTC, []string{"bk-3099", "bk-3100", "bk-3102", "bk-3101"}},
		{"the van's Wednesday", "2027-09-22", "van", van, []string{"bk-3099", "bk-3100"}},
		{"a day without passes", "2027-09-25", "UTC", time.UTC, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Clone(week)
			got := ids(dayBookings(input, tc.day, tc.zone))
			if !slices.Equal(got, tc.want) {
				t.Errorf("dayBookings(week, %q, %s) = %v, want %v", tc.day, tc.zoneName, got, tc.want)
			}
			if !slices.Equal(ids(input), ids(week)) {
				t.Errorf("dayBookings(week, %q, %s) left week as %v, want it unchanged", tc.day, tc.zoneName, ids(input))
			}
		})
	}
}

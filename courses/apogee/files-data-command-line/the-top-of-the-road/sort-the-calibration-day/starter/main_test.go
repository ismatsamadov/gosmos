package main

import (
	"slices"
	"testing"
	"time"
)

// utc returns a time on Thursday 21 October 2027, in UTC.
func utc(hour, min, sec int) time.Time {
	return time.Date(2027, 10, 21, hour, min, sec, 0, time.UTC)
}

// ids returns the bookings' IDs, in order.
func ids(bookings []Booking) []string {
	var out []string
	for _, b := range bookings {
		out = append(out, b.ID)
	}
	return out
}

func TestOrderedPasses(t *testing.T) {
	van := time.FixedZone("", -7*60*60) // the salt flat's clock, UTC-7
	for _, tc := range []struct {
		name   string
		passes []Booking
		want   []string
	}{
		{"out of order", []Booking{{ID: "bk-3417", AOS: utc(17, 22, 40)}, {ID: "bk-3410", AOS: utc(14, 2, 10)}, {ID: "bk-3413", AOS: utc(15, 40, 5)}},
			[]string{"bk-3410", "bk-3413", "bk-3417"}},
		{"one written on the van's clock", []Booking{{ID: "bk-3413", AOS: utc(15, 40, 5)}, {ID: "bk-3419", AOS: time.Date(2027, 10, 21, 8, 30, 0, 0, van)}, {ID: "bk-3410", AOS: utc(14, 2, 10)}},
			[]string{"bk-3410", "bk-3419", "bk-3413"}},
		{"two dishes at 15:40:05", []Booking{{ID: "bk-3414", AOS: utc(15, 40, 5)}, {ID: "bk-3413", AOS: utc(15, 40, 5)}, {ID: "bk-3410", AOS: utc(14, 2, 10)}},
			[]string{"bk-3410", "bk-3413", "bk-3414"}},
		{"two dishes, listed the other way", []Booking{{ID: "bk-3413", AOS: utc(15, 40, 5)}, {ID: "bk-3410", AOS: utc(14, 2, 10)}, {ID: "bk-3414", AOS: utc(15, 40, 5)}},
			[]string{"bk-3410", "bk-3413", "bk-3414"}},
		{"no passes", nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := slices.Clone(tc.passes)
			got := orderedPasses(tc.passes)
			if !slices.Equal(ids(got), tc.want) {
				t.Errorf("orderedPasses(%s) = %q, want %q", tc.name, ids(got), tc.want)
			}
			if !slices.Equal(tc.passes, before) {
				t.Errorf("orderedPasses(%s) changed its input to %q, want %q left as it was", tc.name, ids(tc.passes), ids(before))
			}
		})
	}
}

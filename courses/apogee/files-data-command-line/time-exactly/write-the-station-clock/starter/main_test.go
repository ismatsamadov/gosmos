package main

import (
	"testing"
	"time"
)

func TestHeading(t *testing.T) {
	for _, tc := range []struct {
		name string
		t    time.Time
		want string
	}{
		{"afternoon pass", time.Date(2027, time.September, 21, 14, 2, 10, 0, time.UTC), "Tue 21 Sep 2027, 14:02 UTC"},
		{"seconds left out", time.Date(2027, time.September, 21, 14, 2, 59, 0, time.UTC), "Tue 21 Sep 2027, 14:02 UTC"},
		{"just after midnight", time.Date(2027, time.September, 23, 0, 7, 0, 0, time.UTC), "Thu 23 Sep 2027, 00:07 UTC"},
		{"late evening", time.Date(2027, time.September, 26, 23, 59, 0, 0, time.UTC), "Sun 26 Sep 2027, 23:59 UTC"},
		{"first of the month", time.Date(2027, time.October, 1, 9, 30, 0, 0, time.UTC), "Fri 1 Oct 2027, 09:30 UTC"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := heading(tc.t); got != tc.want {
				t.Errorf("heading(%s) = %q, want %q", tc.t.Format(time.RFC3339), got, tc.want)
			}
		})
	}
}

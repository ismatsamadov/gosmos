package main

import (
	"testing"
	"time"
)

func TestVanTime(t *testing.T) {
	at := time.Date(2027, time.September, 23, 3, 30, 0, 0, time.UTC)
	for _, tc := range []struct {
		name        string
		secondsEast int
		want        string
	}{
		{"the van", -7 * 60 * 60, "2027-09-22T20:30:00-07:00"},
		{"UTC itself", 0, "2027-09-23T03:30:00Z"},
		{"east of UTC", 5*60*60 + 30*60, "2027-09-23T09:00:00+05:30"},
		{"half an hour west", -30 * 60, "2027-09-23T03:00:00-00:30"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := vanTime(at, tc.secondsEast)
			_, offset := got.Zone()
			if got.Format(time.RFC3339) != tc.want || offset != tc.secondsEast {
				t.Errorf("vanTime(%s, %d) = %s with offset %d, want %s with offset %d", at.Format(time.RFC3339), tc.secondsEast, got.Format(time.RFC3339), offset, tc.want, tc.secondsEast)
			}
			if !got.Equal(at) {
				t.Errorf("vanTime(%s, %d) = %s, want the same instant as %s", at.Format(time.RFC3339), tc.secondsEast, got.Format(time.RFC3339), at.Format(time.RFC3339))
			}
		})
	}
}

package main

import (
	"testing"
	"time"
)

func TestParseStart(t *testing.T) {
	want := time.Date(2027, time.September, 21, 14, 2, 10, 0, time.UTC)
	for _, tc := range []struct{ name, text string }{
		{"network clock", "2027-09-21T14:02:10Z"},
		{"van clock", "2027-09-21T07:02:10-07:00"},
		{"east of UTC", "2027-09-21T19:32:10+05:30"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseStart(tc.text)
			if err != nil || !got.Equal(want) {
				t.Errorf("parseStart(%q) = %s, %v, want %s, nil", tc.text, got.Format(time.RFC3339), err, want.Format(time.RFC3339))
			}
		})
	}
	for _, tc := range []struct{ name, text string }{
		{"no offset", "2027-09-21T14:02:10"},
		{"space for T", "2027-09-21 14:02:10Z"},
		{"day first", "21/09/2027 14:02"},
		{"empty", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseStart(tc.text)
			if err == nil || !got.IsZero() {
				t.Errorf("parseStart(%q) = %s, %v, want the zero time and an error", tc.text, got.Format(time.RFC3339), err)
			}
		})
	}
}

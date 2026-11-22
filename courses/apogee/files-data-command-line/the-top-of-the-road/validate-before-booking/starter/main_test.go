package main

import (
	"strings"
	"testing"
	"time"
)

// stations returns the five stations of the station table, by code.
func stations() map[string]Station {
	return map[string]Station{
		"YST": {Code: "YST", Name: "Ystmark"},
		"PEL": {Code: "PEL", Name: "Port Ellery"},
		"SQT": {Code: "SQT", Name: "Salar Quieto"},
		"LAN": {Code: "LAN", Name: "Lantern Atoll"},
		"CWR": {Code: "CWR", Name: "Cape Wren"},
	}
}

// thursday returns a valid batch from Hollis's sheet.
func thursday() []Booking {
	return []Booking{
		{ID: "bk-3410", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2027, 10, 21, 14, 2, 10, 0, time.UTC), LOS: time.Date(2027, 10, 21, 14, 13, 52, 0, time.UTC)},
		{ID: "bk-3411", SatID: "furrow-a", Station: "CWR", AOS: time.Date(2027, 10, 21, 15, 44, 20, 0, time.UTC), LOS: time.Date(2027, 10, 21, 15, 55, 48, 0, time.UTC)},
		{ID: "bk-3422", SatID: "furrow-b", Station: "SQT", AOS: time.Date(2027, 10, 22, 16, 5, 0, 0, time.UTC), LOS: time.Date(2027, 10, 22, 16, 16, 41, 0, time.UTC)},
	}
}

func TestValidateBookingsAccepts(t *testing.T) {
	if err := validateBookings(thursday(), stations()); err != nil {
		t.Errorf("validateBookings(thursday) = %v, want nil", err)
	}
	if err := validateBookings(nil, stations()); err != nil {
		t.Errorf("validateBookings(no bookings) = %v, want nil", err)
	}
}

func TestValidateBookingsRefuses(t *testing.T) {
	for _, tc := range []struct {
		name    string
		spoil   func(b []Booking)
		mention string // the booking the error must name
	}{
		{"no id", func(b []Booking) { b[1].ID = "" }, "2"},
		{"no aos", func(b []Booking) { b[2].AOS = time.Time{} }, "bk-3422"},
		{"no los", func(b []Booking) { b[2].LOS = time.Time{} }, "bk-3422"},
		{"los equal to aos", func(b []Booking) { b[0].LOS = b[0].AOS }, "bk-3410"},
		{"los before aos", func(b []Booking) { b[0].AOS, b[0].LOS = b[0].LOS, b[0].AOS }, "bk-3410"},
		{"sqx on the last line", func(b []Booking) { b[2].Station = "SQX" }, "bk-3422"},
		{"code in lowercase", func(b []Booking) { b[2].Station = "sqt" }, "bk-3422"},
		{"no station", func(b []Booking) { b[1].Station = "" }, "bk-3411"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bookings := thursday()
			tc.spoil(bookings)
			err := validateBookings(bookings, stations())
			if err == nil {
				t.Errorf("validateBookings(%s) = nil, want an error", tc.name)
			} else if !strings.Contains(err.Error(), tc.mention) {
				t.Errorf("validateBookings(%s) = %q, want an error naming %q", tc.name, err, tc.mention)
			}
		})
	}
}

package main

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

const export = `{"id":"bk-2041","sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z","note":"priority"}`

func TestDecodeBooking(t *testing.T) {
	want := Booking{
		ID:      "bk-2041",
		Station: "SQT",
		AOS:     time.Date(2027, 8, 11, 14, 2, 10, 0, time.UTC),
		LOS:     time.Date(2027, 8, 11, 14, 11, 52, 0, time.UTC),
	}
	got, err := decodeBooking([]byte(export))
	if got.ID != want.ID || got.Station != want.Station || !got.AOS.Equal(want.AOS) || !got.LOS.Equal(want.LOS) || err != nil {
		t.Errorf("decodeBooking(bk-2041) = %+v, %v, want %+v, nil", got, err, want)
	}
}

func TestDecodeBookingErrors(t *testing.T) {
	isSyntax := func(err error) bool {
		var target *json.SyntaxError
		return errors.As(err, &target)
	}
	isType := func(err error) bool {
		var target *json.UnmarshalTypeError
		return errors.As(err, &target)
	}
	for _, tc := range []struct {
		name, data, want string
		check            func(error) bool
	}{
		{"cut short", `{"id":"bk-2041","station":`, "a syntax error", isSyntax},
		{"not JSON", `bk-2041 SQT`, "a syntax error", isSyntax},
		{"a number for the id", `{"id":2041,"station":"SQT"}`, "a type error", isType},
		{"a number for the station", `{"id":"bk-2041","station":7}`, "a type error", isType},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeBooking([]byte(tc.data))
			if got != (Booking{}) || !tc.check(err) {
				t.Errorf("decodeBooking(%s) = %+v, %v, want zero Booking, %s", tc.data, got, err, tc.want)
			}
		})
	}
}

package main

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// describe says what checkPass returned: nil, or an *APIError's status, code and message.
func describe(err error) string {
	if err == nil {
		return "nil"
	}
	if e, ok := errors.AsType[*APIError](err); ok {
		return fmt.Sprintf("%d %s %q", e.Status, e.Code, e.Message)
	}
	return fmt.Sprintf("a plain error %q", err)
}

func TestCheckPass(t *testing.T) {
	// Friday 3 December 2027 at 14:30 UTC, 15:30 in Furrowlight's time.
	friday := time.Date(2027, 12, 3, 14, 30, 0, 0, time.UTC)
	ourTime := time.FixedZone("", 60*60)
	hollis := Booking{SatID: "furrow-a", Station: "SQT", AOS: mustTime("2027-12-06T14:02:10Z"), LOS: mustTime("2027-12-06T14:11:52Z")}
	// with returns Hollis's booking changed by change.
	with := func(change func(b *Booking)) Booking {
		b := hollis
		change(&b)
		return b
	}
	tests := []struct {
		name string
		b    Booking
		now  time.Time
		want string // the 422's message, or "" for nil
	}{
		{"Monday's 14:02 over Salar Quieto, booked on Friday", hollis, friday, ""},
		{"the same pass in Furrowlight's time", with(func(b *Booking) { b.AOS, b.LOS = b.AOS.In(ourTime), b.LOS.In(ourTime) }), friday, ""},
		{"Thistle", with(func(b *Booking) { b.SatID = "thistle" }), friday, `unknown satellite "thistle"`},
		{"Thistle at a station nobody has", with(func(b *Booking) { b.SatID, b.Station = "thistle", "ZZZ" }), friday, `unknown satellite "thistle"`},
		{"a station code in lowercase", with(func(b *Booking) { b.Station = "sqt" }), friday, `unknown station "sqt"`},
		{"LOS one second off", with(func(b *Booking) { b.LOS = b.LOS.Add(time.Second) }), friday, "no predicted pass of furrow-a over SQT at 2027-12-06T14:02:10Z"},
		{"Furrow-B at Furrow-A's time", with(func(b *Booking) { b.SatID = "furrow-b" }), friday, "no predicted pass of furrow-b over SQT at 2027-12-06T14:02:10Z"},
		{"no such pass, written in Furrowlight's time", with(func(b *Booking) { b.AOS, b.LOS = b.AOS.In(ourTime), b.LOS.Add(time.Minute).In(ourTime) }), friday, "no predicted pass of furrow-a over SQT at 2027-12-06T14:02:10Z"},
		{"a pass that never existed, last Monday", with(func(b *Booking) { b.AOS, b.LOS = b.AOS.AddDate(0, 0, -7), b.LOS.AddDate(0, 0, -7) }), friday, "no predicted pass of furrow-a over SQT at 2027-11-29T14:02:10Z"},
		{"the pass that began an hour ago", Booking{SatID: "furrow-a", Station: "SQT", AOS: mustTime("2027-12-03T13:28:42Z"), LOS: mustTime("2027-12-03T13:36:51Z")}, friday, "bookings for this pass closed at 2027-12-03T13:18:42Z"},
		{"at the very second bookings close", hollis, mustTime("2027-12-06T13:52:10Z"), "bookings for this pass closed at 2027-12-06T13:52:10Z"},
		{"one second before bookings close", hollis, mustTime("2027-12-06T13:52:09Z"), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPass(tt.b, catalog(), tt.now)
			if tt.want == "" {
				if err != nil {
					t.Errorf("checkPass(%s) = %s, want nil", tt.name, describe(err))
				}
				return
			}
			if e, ok := errors.AsType[*APIError](err); !ok || e.Status != 422 || e.Code != "invalid_booking" || e.Message != tt.want {
				t.Errorf("checkPass(%s) = %s, want 422 invalid_booking %q", tt.name, describe(err), tt.want)
			}
		})
	}
}

package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// describe says what checkShape returned: nil, or an *APIError's status, code and message.
func describe(err error) string {
	if err == nil {
		return "nil"
	}
	if e, ok := errors.AsType[*APIError](err); ok {
		return fmt.Sprintf("%d %s %q", e.Status, e.Code, e.Message)
	}
	return fmt.Sprintf("a plain error %q", err)
}

func TestCheckShape(t *testing.T) {
	aos := time.Date(2027, 12, 6, 14, 2, 10, 0, time.UTC)
	los := time.Date(2027, 12, 6, 14, 11, 52, 0, time.UTC)
	hollis := Booking{SatID: "furrow-a", Station: "SQT", AOS: aos, LOS: los, Note: "PRIORITY"}
	// with returns Hollis's booking changed by change.
	with := func(change func(b *Booking)) Booking {
		b := hollis
		change(&b)
		return b
	}
	tests := []struct {
		name string
		b    Booking
		want string // the 422's message, or "" for nil
	}{
		{"Hollis's booking", hollis, ""},
		{"an ID the caller chose", with(func(b *Booking) { b.ID = "bk-4501" }), "id is chosen by the service"},
		{"no satellite", with(func(b *Booking) { b.SatID = "" }), "sat_id is required"},
		{"no station", with(func(b *Booking) { b.Station = "" }), "station is required"},
		{"no los", with(func(b *Booking) { b.LOS = time.Time{} }), "aos and los are required"},
		{"no aos", with(func(b *Booking) { b.AOS = time.Time{} }), "aos and los are required"},
		{"los at the very second of aos", with(func(b *Booking) { b.LOS = aos }), "los must be after aos"},
		{"los before aos", with(func(b *Booking) { b.AOS, b.LOS = los, aos }), "los must be after aos"},
		{"a note of 200 characters, with Ñ", with(func(b *Booking) { b.Note = "AÑIL " + strings.Repeat("Ñ", 195) }), ""},
		{"a note of 201 characters", with(func(b *Booking) { b.Note = strings.Repeat("A", 201) }), "note is longer than 200 characters"},
		{"no station and no los", with(func(b *Booking) { b.Station, b.LOS = "", time.Time{} }), "station is required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkShape(tt.b)
			if tt.want == "" {
				if err != nil {
					t.Errorf("checkShape(%s) = %s, want nil", tt.name, describe(err))
				}
				return
			}
			if e, ok := errors.AsType[*APIError](err); !ok || e.Status != 422 || e.Code != "invalid_booking" || e.Message != tt.want {
				t.Errorf("checkShape(%s) = %s, want 422 invalid_booking %q", tt.name, describe(err), tt.want)
			}
		})
	}
}

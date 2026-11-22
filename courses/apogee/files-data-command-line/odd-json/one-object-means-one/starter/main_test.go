package main

import (
	"errors"
	"strings"
	"testing"
)

// settingsFile is Paz's settings for tomorrow's test passes.
const settingsFile = `{"station":"SQT","day":"2027-08-26","min_elev_deg":10}`

func TestDecodeOne(t *testing.T) {
	want := Settings{Station: "SQT", Day: "2027-08-26", MinElevDeg: 10}
	for _, tail := range []string{"", "\n", " \n\t\n"} {
		got, err := decodeOne(strings.NewReader(settingsFile + tail))
		if got != want || err != nil {
			t.Errorf("decodeOne(object + %q) = %+v, %v, want %+v, nil", tail, got, err, want)
		}
	}
	for _, tc := range []struct{ name, tail string }{
		{"second object", "\n" + `{"station":"YST","day":"2027-08-26","min_elev_deg":5}`},
		{"object glued on", `{"station":"YST"}`},
		{"stray brace", "}"},
		{"stray word", "\nyes\n"},
		{"null", " null"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeOne(strings.NewReader(settingsFile + tc.tail))
			if !errors.Is(err, ErrTrailing) || got != (Settings{}) {
				t.Errorf("decodeOne(object + %q) = %+v, %v, want empty settings and ErrTrailing", tc.tail, got, err)
			}
		})
	}
	if _, err := decodeOne(strings.NewReader(`{"min_elev_dg":10}`)); err == nil {
		t.Error("decodeOne(misspelled field) error = nil, want an error")
	}
}

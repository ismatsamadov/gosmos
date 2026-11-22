package main

import (
	"strings"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	thirty := Delay(30 * time.Second)
	for _, tc := range []struct {
		name, text string
		want       Settings
	}{
		{"day only", `{"day":"2027-10-21"}`, Settings{Station: "SQT", Day: "2027-10-21", MinElevDeg: 3, Timeout: thirty}},
		{"every field", `{"station":"CWR","day":"2027-10-22","min_elev_deg":5,"timeout":"2m"}`, Settings{Station: "CWR", Day: "2027-10-22", MinElevDeg: 5, Timeout: Delay(2 * time.Minute)}},
		{"explicit zero", `{"day":"2027-10-21","min_elev_deg":0}`, Settings{Station: "SQT", Day: "2027-10-21", MinElevDeg: 0, Timeout: thirty}},
		{"newline after the object", "{\"day\":\"2027-10-21\"}\n", Settings{Station: "SQT", Day: "2027-10-21", MinElevDeg: 3, Timeout: thirty}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadConfig(strings.NewReader(tc.text))
			if got != tc.want || err != nil {
				t.Errorf("loadConfig(%q) = %+v, %v, want %+v, nil", tc.text, got, err, tc.want)
			}
		})
	}
}

func TestLoadConfigRefuses(t *testing.T) {
	for _, tc := range []struct{ name, text string }{
		{"misspelled field", `{"day":"2027-10-21","min_elev":5}`},
		{"second object", `{"day":"2027-10-21"} {"day":"2027-10-22"}`},
		{"text after the object", `{"day":"2027-10-21"} thursday`},
		{"no day", `{"station":"SQT"}`},
		{"day that is not a date", `{"day":"../outbox"}`},
		{"timeout as a number", `{"day":"2027-10-21","timeout":30}`},
		{"empty file", ``},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadConfig(strings.NewReader(tc.text))
			if err == nil {
				t.Errorf("loadConfig(%q) = %+v, nil, want an error", tc.text, got)
			}
		})
	}
}

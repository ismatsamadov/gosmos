package main

import (
	"strings"
	"testing"
)

func TestStrictConfig(t *testing.T) {
	for _, tc := range []struct {
		name, file string
		want       Settings
	}{
		{"every field", `{"station":"SQT","day":"2027-08-26","min_elev_deg":10}`, Settings{Station: "SQT", Day: "2027-08-26", MinElevDeg: 10}},
		{"day only", `{"day":"2027-08-27"}`, Settings{Day: "2027-08-27"}},
		{"spread over lines", "{\n  \"min_elev_deg\": 3.5,\n  \"station\": \"SQT\"\n}\n", Settings{Station: "SQT", MinElevDeg: 3.5}},
	} {
		got, err := strictConfig(strings.NewReader(tc.file))
		if got != tc.want || err != nil {
			t.Errorf("strictConfig(%s) = %+v, %v, want %+v, nil", tc.name, got, err, tc.want)
		}
	}
	for _, tc := range []struct{ name, file, field string }{
		{"misspelled limit", `{"station":"SQT","day":"2027-08-26","min_elev_dg":10}`, "min_elev_dg"},
		{"extra field", `{"station":"SQT","day":"2027-08-26","dish":"LG"}`, "dish"},
	} {
		got, err := strictConfig(strings.NewReader(tc.file))
		if err == nil || !strings.Contains(err.Error(), tc.field) || got != (Settings{}) {
			t.Errorf("strictConfig(%s) = %+v, %v, want empty settings and an error naming %q", tc.name, got, err, tc.field)
		}
	}
	if _, err := strictConfig(strings.NewReader(`{"station":"SQT",`)); err == nil {
		t.Error("strictConfig(cut off) error = nil, want an error")
	}
}

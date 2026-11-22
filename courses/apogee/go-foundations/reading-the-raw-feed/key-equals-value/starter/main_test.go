package main

import "testing"

// thistleLine is a telemetry line with a note whose value holds an "=" of its own.
const thistleLine = "thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL note=a=b"

func TestField(t *testing.T) {
	cases := []struct {
		name, key, want string
		wantOK          bool
	}{
		{"battery", "batt_v", "7.91", true},
		{"mode", "mode", "NOMINAL", true},
		{"missing key", "wind_kmh", "", false},
		{"value with =", "note", "a=b", true},
		{"satellite is no key", "thistle", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := field(thistleLine, c.key)
			if got != c.want || ok != c.wantOK {
				t.Errorf("field(line, %q) = %q, %v, want %q, %v", c.key, got, ok, c.want, c.wantOK)
			}
		})
	}
}

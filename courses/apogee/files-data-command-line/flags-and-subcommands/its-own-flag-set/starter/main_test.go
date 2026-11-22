package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseOptions(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"every station", nil, ""},
		{"Salar Quieto", []string{"-station", "SQT"}, "SQT"},
		{"with an equals sign", []string{"-station=CWR"}, "CWR"},
		{"with two dashes", []string{"--station", "PEL"}, "PEL"},
		{"before a file name", []string{"-station", "LAN", "week.json"}, "LAN"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			got, err := parseOptions(tc.args, &output)
			if got.Station != tc.want || err != nil {
				t.Errorf("parseOptions(%q) = %+v, %v, want {Station:%s}, nil", tc.args, got, err, tc.want)
			}
			if output.Len() != 0 {
				t.Errorf("parseOptions(%q) wrote %q, want nothing", tc.args, output.String())
			}
		})
	}
	for _, tc := range []struct {
		name string
		args []string
		flag string
	}{
		{"misspelled flag", []string{"-stn", "SQT"}, "-stn"},
		{"no value", []string{"-station"}, "-station"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			got, err := parseOptions(tc.args, &output)
			if got != (passesOptions{}) || err == nil {
				t.Errorf("parseOptions(%q) = %+v, %v, want empty options and an error", tc.args, got, err)
			}
			if !strings.Contains(output.String(), tc.flag) {
				t.Errorf("parseOptions(%q) wrote %q to output, want a message naming %s", tc.args, output.String(), tc.flag)
			}
		})
	}
}

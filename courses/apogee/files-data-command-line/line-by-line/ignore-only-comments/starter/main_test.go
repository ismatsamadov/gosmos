package main

import (
	"slices"
	"strings"
	"testing"
)

const (
	note = "# night shift 2027-07-06, Paz"
	lg   = "2027-07-06T03:14:07Z\tLG\t183.25\t41.70"
	ch   = "2027-07-06T03:14:08Z  CH  183.31   41.66"
)

func TestDataLines(t *testing.T) {
	for _, tc := range []struct {
		name, log string
		want      []string
	}{
		{"a note and a blank line", note + "\n" + lg + "\n\n" + ch + "\n", []string{lg, ch}},
		{"indented notes", "\t# LG warm after noon, Paz\n  # wind 40 km/h\n" + lg + "\n", []string{lg}},
		{"spaces and tabs only", "   \n\t\n \t \n" + ch + "\n", []string{ch}},
		{"indented readings", "  " + ch + "\n\t" + lg + "\n", []string{"  " + ch, "\t" + lg}},
		{"notes only", note + "\n\n", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := dataLines(strings.NewReader(tc.log))
			if !slices.Equal(got, tc.want) || err != nil {
				t.Errorf("dataLines(%s) = %q, %v, want %q, nil", tc.name, got, err, tc.want)
			}
		})
	}
}

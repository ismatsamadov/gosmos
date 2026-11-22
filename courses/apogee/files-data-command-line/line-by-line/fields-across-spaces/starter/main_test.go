package main

import (
	"slices"
	"testing"
)

func TestParseColumns(t *testing.T) {
	want := []string{"2027-07-06T03:14:07Z", "LG", "183.25", "41.70"}
	for _, tc := range []struct{ name, line string }{
		{"tabs", "2027-07-06T03:14:07Z\tLG\t183.25\t41.70"},
		{"repeated spaces", "2027-07-06T03:14:07Z  LG  183.25   41.70"},
		{"tabs and spaces together", "2027-07-06T03:14:07Z \tLG\t  183.25 \t41.70"},
		{"indented", "  2027-07-06T03:14:07Z\tLG\t183.25\t41.70\t"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseColumns(tc.line)
			if !slices.Equal(got, want) || err != nil {
				t.Errorf("parseColumns(%q) = %q, %v, want %q, nil", tc.line, got, err, want)
			}
		})
	}
	for _, tc := range []struct {
		name, line, want string
	}{
		{"missing elevation", "2027-07-06T03:14:07Z\tLG\t183.25", "3 columns instead of 4"},
		{"one column too many", "2027-07-06T03:14:08Z  CH  183.31  41.66  41.70", "5 columns instead of 4"},
		{"blank line", "", "0 columns instead of 4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseColumns(tc.line)
			if got != nil || err == nil || err.Error() != tc.want {
				t.Errorf("parseColumns(%q) = %q, %v, want nil, %s", tc.line, got, err, tc.want)
			}
		})
	}
}

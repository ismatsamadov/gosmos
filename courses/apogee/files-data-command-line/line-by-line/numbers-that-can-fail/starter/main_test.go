package main

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

func TestParseElevation(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		want       float64
	}{
		{"La Grande tracking", "41.70", 41.70},
		{"straight up", "90", 90},
		{"flat on the horizon", "0", 0},
		{"a little below the horizon", "-0.35", -0.35},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseElevation(tc.text)
			if math.Abs(got-tc.want) > 0.01 || err != nil {
				t.Errorf("parseElevation(%q) = %.2f, %v, want %.2f, nil", tc.text, got, err, tc.want)
			}
		})
	}
	for _, tc := range []struct{ name, text string }{
		{"a doubled decimal point", "41..70"},
		{"a comma for a point", "41,70"},
		{"an empty column", ""},
		{"a word", "forty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseElevation(tc.text)
			if got != 0 || !errors.Is(err, strconv.ErrSyntax) {
				t.Errorf("parseElevation(%q) = %.2f, %v, want 0.00, %v", tc.text, got, err, strconv.ErrSyntax)
			}
		})
	}
}

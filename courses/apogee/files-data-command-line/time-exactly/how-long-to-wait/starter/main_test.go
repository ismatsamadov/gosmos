package main

import (
	"testing"
	"time"
)

func TestPositiveDelay(t *testing.T) {
	for _, tc := range []struct {
		text string
		want time.Duration
	}{
		{"30s", 30 * time.Second},
		{"2m", 2 * time.Minute},
		{"1m30s", 90 * time.Second},
		{"250ms", 250 * time.Millisecond},
	} {
		t.Run(tc.text, func(t *testing.T) {
			got, err := positiveDelay(tc.text)
			if got != tc.want || err != nil {
				t.Errorf("positiveDelay(%q) = %v, %v, want %v, nil", tc.text, got, err, tc.want)
			}
		})
	}
	for _, tc := range []struct{ name, text string }{
		{"zero seconds", "0s"},
		{"bare zero", "0"},
		{"negative", "-5s"},
		{"no unit", "30"},
		{"words", "two minutes"},
		{"empty", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := positiveDelay(tc.text)
			if got != 0 || err == nil {
				t.Errorf("positiveDelay(%q) = %v, %v, want 0 and an error", tc.text, got, err)
			}
		})
	}
}

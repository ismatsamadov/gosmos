package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestCheckTemp(t *testing.T) {
	tests := []struct {
		name  string
		c     float64
		limit bool // whether checkTemp(c) returns a *LimitError
	}{
		{"thistle at 51 C", 51, true},
		{"thistle at -12.5 C", -12.5, true},
		{"pipit-1 at 45 C", 45, false},
		{"pipit-2 at 20 C", 20, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkTemp(tt.c)
			if !tt.limit {
				if err != nil {
					t.Errorf("checkTemp(%v) = %v, want nil", tt.c, err)
				}
				return
			}
			le, ok := errors.AsType[*LimitError](err)
			if !ok {
				t.Fatalf("checkTemp(%v) = %v, want a *LimitError", tt.c, err)
			}
			want := LimitError{Field: "temp_c", Value: tt.c, Min: -10, Max: 45}
			if *le != want {
				t.Errorf("checkTemp(%v) = %+v, want %+v", tt.c, *le, want)
			}
		})
	}
}

func TestPagerDetail(t *testing.T) {
	hot := &LimitError{Field: "temp_c", Value: 51, Min: -10, Max: 45}
	cold := &LimitError{Field: "temp_c", Value: -12.5, Min: -10, Max: 45}
	flat := &LimitError{Field: "batt_v", Value: 6.2, Min: 6.6, Max: 8.4}
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"wrapped", fmt.Errorf("thistle seq 4190: %w", hot), "temp_c=51.0 outside [-10.0, 45.0]"},
		{"bare", cold, "temp_c=-12.5 outside [-10.0, 45.0]"},
		{"battery", fmt.Errorf("thistle seq 4191: %w", flat), "batt_v=6.2 outside [6.6, 8.4]"},
		{"bad crc", errors.New("bad crc"), ""},
		{"nil", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pagerDetail(tt.err); got != tt.want {
				t.Errorf("pagerDetail(%v) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}

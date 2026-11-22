package main

import (
	"testing"
	"time"
)

// nine is 09:00:00 UTC on Monday 10 January, when each test's bucket was last brought up to date.
var nine = time.Date(2028, 1, 10, 9, 0, 0, 0, time.UTC)

// requests is a burst of n requests at a moment after nine, and how many of them Take must let through.
type requests struct {
	after time.Duration
	n     int
	want  int
}

func TestTake(t *testing.T) {
	full := Bucket{Tokens: 10, Last: nine}
	empty := Bucket{Tokens: 0, Last: nine}
	tests := []struct {
		name   string
		bucket Bucket
		bursts []requests
	}{
		{"ten at once, then no more", full, []requests{{0, 11, 10}}},
		{"one more a second later", empty, []requests{{time.Second, 2, 1}}},
		{"half a second isn't a whole token", empty, []requests{{500 * time.Millisecond, 1, 0}, {time.Second, 2, 1}}},
		{"an hour later, ten, not 3,600", empty, []requests{{time.Hour, 11, 10}}},
		{"a clock stepping backward adds nothing", empty, []requests{{-10 * time.Second, 1, 0}, {time.Second, 2, 1}}},
		{"a zero Bucket starts full", Bucket{}, []requests{{0, 11, 10}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.bucket
			for _, burst := range tt.bursts {
				at := nine.Add(burst.after)
				got := 0
				for range burst.n {
					if b.Take(at) {
						got++
					}
				}
				if got != burst.want {
					t.Errorf("Take(%s) for a burst of %d = %d allowed, want %d", at.Format("15:04:05.000"), burst.n, got, burst.want)
				}
			}
		})
	}
}

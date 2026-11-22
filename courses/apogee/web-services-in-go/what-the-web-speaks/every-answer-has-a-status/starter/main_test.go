package main

import "testing"

func TestAdvice(t *testing.T) {
	tests := []struct {
		name string
		code int
		want string
	}{
		{"the station list", 200, "done"},
		{"a new booking", 201, "done"},
		{"a cancellation", 204, "done"},
		{"the last of the 2xx", 299, "done"},
		{"a broken request", 400, "fix the request"},
		{"no station ZZZ", 404, "fix the request"},
		{"a pass already booked", 409, "fix the request"},
		{"the last of the 4xx", 499, "fix the request"},
		{"too many requests too soon", 429, "wait and retry"},
		{"a bug in the service", 500, "retry later"},
		{"the service restarting", 503, "retry later"},
		{"the last of the 5xx", 599, "retry later"},
		{"one below the 2xx", 199, "unexpected"},
		{"a redirect", 302, "unexpected"},
		{"one past the 5xx", 600, "unexpected"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := advice(tt.code); got != tt.want {
				t.Errorf("advice(%d) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

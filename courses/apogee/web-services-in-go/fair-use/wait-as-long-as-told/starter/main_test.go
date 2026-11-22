package main

import (
	"net/http"
	"testing"
	"time"
)

func TestRetryDelay(t *testing.T) {
	tests := []struct {
		name  string
		value string // "" sends no Retry-After header
		want  time.Duration
	}{
		{"the service's usual answer", "1", time.Second},
		{"three seconds", "3", 3 * time.Second},
		{"a whole minute", "60", time.Minute},
		{"one second too many", "61", time.Minute},
		{"ten minutes", "600", time.Minute},
		{"more seconds than a Duration holds", "9999999999", time.Minute},
		{"no header", "", time.Second},
		{"a word", "soon", time.Second},
		{"half seconds", "2.5", time.Second},
		{"zero", "0", time.Second},
		{"below zero", "-5", time.Second},
		{"a date, which HTTP allows", "Fri, 14 Jan 2028 09:00:05 GMT", time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{}}
			call := "retryDelay(no Retry-After)"
			if tt.value != "" {
				res.Header.Set("Retry-After", tt.value)
				call = "retryDelay(Retry-After: " + tt.value + ")"
			}
			if got := retryDelay(res); got != tt.want {
				t.Errorf("%s = %v, want %v", call, got, tt.want)
			}
		})
	}
}

package main

import (
	"bufio"
	"errors"
	"strings"
	"testing"
)

const reading = "2027-07-06T03:14:07Z\tLG\t183.25\t41.70\n"

// sweep returns one calibration sweep n bytes long, newline included: a time, a dish, then every sample.
func sweep(n int) string {
	line := "2027-07-02T15:00:00Z\tLG\t" + strings.Repeat("183.250 41.700 ", n/15)
	return line[:n-1] + "\n"
}

func TestScanRecords(t *testing.T) {
	for _, tc := range []struct {
		name, log string
		want      int
	}{
		{"readings only", reading + reading + reading, 3},
		{"a 100 KiB sweep", reading + sweep(100*1024) + reading, 3},
		{"a 200 KiB sweep", reading + sweep(200*1024) + reading, 3},
		{"a sweep of exactly 256 KiB, newline included", reading + sweep(256*1024) + reading, 3},
		{"a sweep without a final newline", reading + strings.TrimSuffix(sweep(200*1024), "\n"), 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := scanRecords(strings.NewReader(tc.log))
			if got != tc.want || err != nil {
				t.Errorf("scanRecords(%s) = %d, %v, want %d, nil", tc.name, got, err, tc.want)
			}
		})
	}
	got, err := scanRecords(strings.NewReader(reading + sweep(300*1024) + reading))
	if got != 1 || !errors.Is(err, bufio.ErrTooLong) {
		t.Errorf("scanRecords(a 300 KiB record) = %d, %v, want 1, %v", got, err, bufio.ErrTooLong)
	}
}

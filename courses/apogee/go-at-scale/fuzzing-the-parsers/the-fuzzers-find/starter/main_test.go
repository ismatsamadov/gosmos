package main

import (
	"testing"
	"time"
)

// safely returns parse wrapped so that a panic fails the test as parseLine("...") panicked: <value>, instead of ending the test binary.
func safely(t *testing.T, parse func(string) (TelemetryFrame, error)) func(string) (TelemetryFrame, error) {
	return func(line string) (TelemetryFrame, error) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parseLine(%q) panicked: %v", line, r)
			}
		}()
		return parse(line)
	}
}

// FuzzParseLine checks that parseLine never panics: on a good line from the drill, and on every
// file saved under testdata/fuzz/FuzzParseLine.
func FuzzParseLine(f *testing.F) {
	f.Add("cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL")
	f.Fuzz(func(t *testing.T, line string) {
		safely(t, parseLine)(line)
	})
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name, line, want string
	}{
		{"Saturday's 05:10 line from Thistle-2, four fields", "thistle-2 2290 2029-04-14T05:10:41Z batt_v=7.40", "line: want 6 fields, got 4"},
		{"an empty line", "", "line: want 6 fields, got 0"},
		{"seven fields", "cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL rssi=-97", "line: want 6 fields, got 7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := safely(t, parseLine)(tt.line); err == nil || err.Error() != tt.want {
				t.Errorf("parseLine(%q) error = %v, want %q", tt.line, err, tt.want)
			}
		})
	}
	t.Run("a good line from the drill", func(t *testing.T) {
		line := "cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL"
		f, err := safely(t, parseLine)(line)
		if err != nil {
			t.Fatalf("parseLine(%q) error = %v, want nil", line, err)
		}
		at := time.Date(2029, 3, 7, 20, 45, 5, 0, time.UTC)
		if f.SatID != "cirrine-1" || f.Seq != 9012 || !f.At.Equal(at) || f.BattV != 8.02 || f.TempC != 11.5 || f.Mode != "NOMINAL" {
			t.Errorf("parseLine(%q) = %+v, want cirrine-1's frame 9012", line, f)
		}
	})
}

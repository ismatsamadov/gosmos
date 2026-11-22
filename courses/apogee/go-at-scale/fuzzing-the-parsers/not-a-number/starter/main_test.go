package main

import "testing"

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

// FuzzParseLine runs checkLine on parseLine for every file under testdata/fuzz/FuzzParseLine:
// Thistle-2's lines of Saturday 14 April, as Joss forwarded them.
func FuzzParseLine(f *testing.F) {
	f.Fuzz(func(t *testing.T, line string) {
		checkLine(t, safely(t, parseLine), line)
	})
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name, line string
		refused    bool
	}{
		{"Saturday's 03:36 line, batt_v=nan", "thistle-2 2214 2029-04-14T03:36:12Z batt_v=nan temp_c=1.0 mode=SAFE", true},
		{"a temperature sensor that reads nan", "thistle-2 2215 2029-04-14T03:37:27Z batt_v=7.42 temp_c=nan mode=SAFE", true},
		{"Saturday's 03:38 line, temp_c=-Inf", "thistle-2 2216 2029-04-14T03:38:42Z batt_v=7.42 temp_c=-Inf mode=SAFE", true},
		{"a flat battery, 0.00 V", "thistle-2 2217 2029-04-14T03:39:57Z batt_v=0.00 temp_c=1.1 mode=SAFE", false},
		{"20.00 V, the top of the range", "thistle-2 2217 2029-04-14T03:39:57Z batt_v=20.00 temp_c=1.1 mode=SAFE", false},
		{"20.01 V, a hundredth past it", "thistle-2 2217 2029-04-14T03:39:57Z batt_v=20.01 temp_c=1.1 mode=SAFE", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := safely(t, parseLine)(tt.line)
			switch {
			case tt.refused && err == nil:
				t.Errorf("parseLine(%q) = batt_v %.2f, temp_c %.1f, want an error", tt.line, f.BattV, f.TempC)
			case !tt.refused && err != nil:
				t.Errorf("parseLine(%q) error = %v, want nil", tt.line, err)
			}
		})
	}
}

package main

import (
	"math"
	"testing"
)

// checkLine fails t when parse accepts line but gives a frame no satellite could send.
func checkLine(t *testing.T, parse func(string) (TelemetryFrame, error), line string) {
	f, err := parse(line)
	if err != nil {
		return // a refused line needs no more checks
	}
	if f.SatID == "" {
		t.Errorf("parseLine(%q).SatID = %q, want a satellite ID", line, f.SatID)
	}
	if math.IsNaN(f.BattV) || f.BattV < 0 || f.BattV > 20 {
		t.Errorf("parseLine(%q).BattV = %.2f, want 0.00 to 20.00", line, f.BattV)
	}
	if math.IsNaN(f.TempC) || f.TempC < -100 || f.TempC > 150 {
		t.Errorf("parseLine(%q).TempC = %.1f, want -100.0 to 150.0", line, f.TempC)
	}
	switch f.Mode {
	case "NOMINAL", "SAFE", "DETUMBLE", "LEOP":
	default:
		t.Errorf("parseLine(%q).Mode = %q, want NOMINAL, SAFE, DETUMBLE or LEOP", line, f.Mode)
	}
}

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// formatLine writes f as a telemetry line, the one form tlmd accepts.
func formatLine(f TelemetryFrame) string {
	return fmt.Sprintf("%s %d %s batt_v=%.2f temp_c=%.1f mode=%s",
		f.SatID, f.Seq, f.At.UTC().Format("2006-01-02T15:04:05Z"), f.BattV, f.TempC, f.Mode)
}

// parseLine reads a telemetry line strictly: it accepts exactly what formatLine writes, and
// otherwise names the first problem it finds, as "line: want 6 fields, got 5".
func parseLine(line string) (TelemetryFrame, error) {
	f := strings.Fields(line)
	if len(f) != 6 {
		return TelemetryFrame{}, fmt.Errorf("line: want 6 fields, got %d", len(f))
	}
	var fr TelemetryFrame
	if !validSatID(f[0]) {
		return TelemetryFrame{}, fmt.Errorf("line: satellite %q: want lowercase letters, digits and -", f[0])
	}
	fr.SatID = f[0]
	seq, err := strconv.ParseUint(f[1], 10, 32)
	if err != nil || f[1] != strconv.FormatUint(seq, 10) {
		return TelemetryFrame{}, fmt.Errorf("line: sequence %q: want a whole number from 0 to 4294967295", f[1])
	}
	fr.Seq = uint32(seq)
	at, err := time.Parse("2006-01-02T15:04:05Z", f[2])
	if err != nil {
		return TelemetryFrame{}, fmt.Errorf("line: time %q: want 2006-01-02T15:04:05Z", f[2])
	}
	fr.At = at
	if fr.BattV, err = reading(f[3], 4, "batt_v", 2, 0, 20, "volts like 7.91"); err != nil {
		return TelemetryFrame{}, err
	}
	if fr.TempC, err = reading(f[4], 5, "temp_c", 1, -100, 150, "degrees like -4.5"); err != nil {
		return TelemetryFrame{}, err
	}
	mode, ok := strings.CutPrefix(f[5], "mode=")
	if !ok {
		return TelemetryFrame{}, fmt.Errorf("line: field 6: want mode=")
	}
	switch Mode(mode) {
	case "NOMINAL", "SAFE", "DETUMBLE", "LEOP":
		fr.Mode = Mode(mode)
	default:
		return TelemetryFrame{}, fmt.Errorf("line: mode %q: want NOMINAL, SAFE, DETUMBLE or LEOP", mode)
	}
	if formatLine(fr) != line {
		return TelemetryFrame{}, fmt.Errorf("line: not canonical")
	}
	return fr, nil
}

// validSatID reports whether id is 1 to 24 lowercase letters, digits and hyphens,
// starting with a letter.
func validSatID(id string) bool {
	if len(id) == 0 || len(id) > 24 || id[0] < 'a' || id[0] > 'z' {
		return false
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}

// reading reads field n, "<key>=<number>", whose number must have exactly decimals digits
// after its point and lie from lo to hi.
func reading(field string, n int, key string, decimals int, lo, hi float64, like string) (float64, error) {
	text, ok := strings.CutPrefix(field, key+"=")
	if !ok {
		return 0, fmt.Errorf("line: field %d: want %s=", n, key)
	}
	v, err := strconv.ParseFloat(text, 64)
	dot := strings.IndexByte(text, '.')
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || dot < 0 || len(text)-dot-1 != decimals {
		return 0, fmt.Errorf("line: %s %q: want %s", key, text, like)
	}
	if v < lo || v > hi {
		return 0, fmt.Errorf("line: %s %s: want %.*f to %.*f", key, text, decimals, lo, decimals, hi)
	}
	return v, nil
}

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseLine reads a telemetry line into a frame, field by field, refusing any line formatLine wouldn't write.
func parseLine(line string) (TelemetryFrame, error) {
	f := strings.Fields(line)
	if len(f) != 6 {
		return TelemetryFrame{}, fmt.Errorf("line: want 6 fields, got %d", len(f))
	}
	if !satelliteID(f[0]) {
		return TelemetryFrame{}, fmt.Errorf("line: satellite %q: want lowercase letters, digits and -", f[0])
	}
	seq, err := strconv.ParseUint(f[1], 10, 32)
	if err != nil {
		return TelemetryFrame{}, fmt.Errorf("line: sequence %q: want a whole number from 0 to 4294967295", f[1])
	}
	at, err := time.Parse("2006-01-02T15:04:05Z", f[2])
	if err != nil {
		return TelemetryFrame{}, fmt.Errorf("line: time %q: want 2006-01-02T15:04:05Z", f[2])
	}
	battV, err := volts(f[3])
	if err != nil {
		return TelemetryFrame{}, err
	}
	tempC, err := degrees(f[4])
	if err != nil {
		return TelemetryFrame{}, err
	}
	mode, ok := strings.CutPrefix(f[5], "mode=")
	if !ok {
		return TelemetryFrame{}, fmt.Errorf("line: field 6: want mode=")
	}
	switch mode {
	case "NOMINAL", "SAFE", "DETUMBLE", "LEOP":
	default:
		return TelemetryFrame{}, fmt.Errorf("line: mode %q: want NOMINAL, SAFE, DETUMBLE or LEOP", mode)
	}
	fr := TelemetryFrame{SatID: f[0], Seq: uint32(seq), At: at, BattV: battV, TempC: tempC, Mode: Mode(mode)}
	if formatLine(fr) != line {
		return TelemetryFrame{}, fmt.Errorf("line: not canonical")
	}
	return fr, nil
}

// satelliteID reports whether id is 1 to 24 lowercase letters, digits and -, starting with a letter.
func satelliteID(id string) bool {
	if len(id) == 0 || len(id) > 24 || id[0] < 'a' || id[0] > 'z' {
		return false
	}
	for _, r := range id {
		lower, digit := 'a' <= r && r <= 'z', '0' <= r && r <= '9'
		if !lower && !digit && r != '-' {
			return false
		}
	}
	return true
}

// volts reads field 4, such as "batt_v=7.91": a battery voltage from 0.00 to 20.00.
func volts(field string) (float64, error) {
	text, ok := strings.CutPrefix(field, "batt_v=")
	if !ok {
		return 0, fmt.Errorf("line: field 4: want batt_v=")
	}
	v, ok := number(text)
	if !ok {
		return 0, fmt.Errorf("line: batt_v %q: want volts like 7.91", text)
	}
	if !(v >= 0 && v <= 20) {
		return 0, fmt.Errorf("line: batt_v %.2f: want 0.00 to 20.00", v)
	}
	return v, nil
}

// degrees reads field 5, such as "temp_c=-4.5": a temperature from -100.0 to 150.0.
func degrees(field string) (float64, error) {
	text, ok := strings.CutPrefix(field, "temp_c=")
	if !ok {
		return 0, fmt.Errorf("line: field 5: want temp_c=")
	}
	v, ok := number(text)
	if !ok {
		return 0, fmt.Errorf("line: temp_c %q: want degrees like -4.5", text)
	}
	if !(v >= -100 && v <= 150) {
		return 0, fmt.Errorf("line: temp_c %.1f: want -100.0 to 150.0", v)
	}
	return v, nil
}

// number reads text written as formatLine writes a reading, with digits, a point and a sign only, so never nan, Inf or 1e1.
func number(text string) (float64, bool) {
	if strings.Trim(text, "+-.0123456789") != "" {
		return 0, false // a character no reading has
	}
	v, err := strconv.ParseFloat(text, 64)
	return v, err == nil
}

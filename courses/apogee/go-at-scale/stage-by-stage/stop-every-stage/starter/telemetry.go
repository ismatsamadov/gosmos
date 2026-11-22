package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Mode is a satellite's operating mode.
type Mode string

// TelemetryFrame is one decoded housekeeping frame.
type TelemetryFrame struct {
	SatID string
	Seq   uint32
	At    time.Time
	BattV float64
	TempC float64
	Mode  Mode
}

// Heard is one telemetry line as a station heard it.
type Heard struct {
	Station string // "TIR"
	Line    string // "cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL"
}

// parseLine reads a telemetry line into a frame, field by field, as tlmd has since 2027.
func parseLine(line string) (TelemetryFrame, error) {
	f := strings.Fields(line)
	seq, err := strconv.ParseUint(f[1], 10, 32)
	if err != nil {
		return TelemetryFrame{}, fmt.Errorf("line: sequence: %w", err)
	}
	at, err := time.Parse(time.RFC3339, f[2])
	if err != nil {
		return TelemetryFrame{}, fmt.Errorf("line: time: %w", err)
	}
	battV, err := reading(f[3], "batt_v=", 0, 20)
	if err != nil {
		return TelemetryFrame{}, err
	}
	tempC, err := reading(f[4], "temp_c=", -100, 150)
	if err != nil {
		return TelemetryFrame{}, err
	}
	mode, ok := strings.CutPrefix(f[5], "mode=")
	if !ok {
		return TelemetryFrame{}, fmt.Errorf("line: field 6: want mode=")
	}
	return TelemetryFrame{SatID: f[0], Seq: uint32(seq), At: at, BattV: battV, TempC: tempC, Mode: Mode(mode)}, nil
}

// reading reads the number in a field such as "batt_v=8.02", refusing one outside lo to hi.
func reading(field, key string, lo, hi float64) (float64, error) {
	text, ok := strings.CutPrefix(field, key)
	if !ok {
		return 0, fmt.Errorf("line: want %s in %q", key, field)
	}
	v, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("line: %s%s: %w", key, text, err)
	}
	if v < lo || v > hi {
		return 0, fmt.Errorf("line: %s%s: want %g to %g", key, text, lo, hi)
	}
	return v, nil
}

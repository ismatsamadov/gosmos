package main

import "time"

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

// Health is each satellite's newest frame (its highest Seq), by satellite ID.
type Health map[string]TelemetryFrame

// Package telemetry holds decoded housekeeping frames.
package telemetry

import "apogee/groundctl/severity"

// Frame is one decoded housekeeping frame.
type Frame struct {
	SatID string
	Seq   uint32
	BattV float64
}

// Grade ranks a frame by its battery: under 7.2 V is a warning, under 6.6 V critical.
func Grade(f Frame) severity.Severity {
	switch {
	case f.BattV < 6.6:
		return severity.Critical
	case f.BattV < 7.2:
		return severity.Warning
	}
	return severity.Info
}

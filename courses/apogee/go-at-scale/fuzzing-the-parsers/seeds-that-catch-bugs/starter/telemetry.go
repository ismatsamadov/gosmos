package main

import (
	"fmt"
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

// formatLine writes f as a telemetry line, the one form tlmd accepts.
func formatLine(f TelemetryFrame) string {
	return fmt.Sprintf("%s %d %s batt_v=%.2f temp_c=%.1f mode=%s",
		f.SatID, f.Seq, f.At.UTC().Format("2006-01-02T15:04:05Z"), f.BattV, f.TempC, f.Mode)
}

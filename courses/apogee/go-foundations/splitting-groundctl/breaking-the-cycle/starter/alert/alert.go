// Package alert writes the lines the pager sends.
package alert

import (
	"fmt"

	"apogee/groundctl/telemetry"
)

// Severity ranks an anomaly.
type Severity int

// The severities, from least to most urgent.
const (
	Info Severity = iota
	Warning
	Critical
)

// String returns the severity as the pager prints it.
func (s Severity) String() string {
	switch s {
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	}
	return "UNKNOWN"
}

// Line describes a frame and its severity for the pager.
func Line(f telemetry.Frame, s Severity) string {
	return fmt.Sprintf("%s %d: %v", f.SatID, f.Seq, s)
}

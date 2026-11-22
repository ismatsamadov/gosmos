// Package alert writes the lines the pager sends.
package alert

import (
	"fmt"

	"apogee/groundctl/severity"
	"apogee/groundctl/telemetry"
)

// Line describes a frame and its severity for the pager.
func Line(f telemetry.Frame, s severity.Severity) string {
	return fmt.Sprintf("%s %d: %v", f.SatID, f.Seq, s)
}

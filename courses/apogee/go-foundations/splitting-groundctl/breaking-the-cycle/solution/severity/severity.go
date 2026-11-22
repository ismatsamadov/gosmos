// Package severity ranks anomalies for every package that needs to.
package severity

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

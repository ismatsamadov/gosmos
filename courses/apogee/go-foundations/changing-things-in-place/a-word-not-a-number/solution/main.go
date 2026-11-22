package main

import "fmt"

// Severity ranks an anomaly.
type Severity int

const (
	Info Severity = iota
	Warning
	Critical
)

// String returns the severity's name, as pages show it.
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

func main() {
	fmt.Println(Info)
	fmt.Println(Warning)
	fmt.Println(Critical)
}

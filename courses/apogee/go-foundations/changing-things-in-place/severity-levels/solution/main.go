package main

import "fmt"

// Severity ranks an anomaly.
type Severity int

const (
	Info Severity = iota
	Warning
	Critical
)

// worst returns the more severe of a and b.
func worst(a, b Severity) Severity {
	if b > a {
		return b
	}
	return a
}

func main() {
	fmt.Println(worst(Warning, Critical))
}

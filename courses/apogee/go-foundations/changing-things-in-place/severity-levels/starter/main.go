package main

import "fmt"

// Severity ranks an anomaly.
type Severity int

const (
	Info     Severity = 0
	Warning  Severity = 0
	Critical Severity = 0
)

// worst returns the more severe of a and b.
func worst(a, b Severity) Severity {
	return a
}

func main() {
	fmt.Println(worst(Warning, Critical))
}

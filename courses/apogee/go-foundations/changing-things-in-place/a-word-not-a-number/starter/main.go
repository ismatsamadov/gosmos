package main

import "fmt"

// Severity ranks an anomaly.
type Severity int

const (
	Info Severity = iota
	Warning
	Critical
)

// TODO: add a String method that names each severity.

func main() {
	fmt.Println(Info)
	fmt.Println(Warning)
	fmt.Println(Critical)
}

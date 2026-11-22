package main

import (
	"fmt"
	"time"
)

// Severity ranks an anomaly.
type Severity int

const (
	Info Severity = iota
	Warning
	Critical
)

// Rule decides when a station's metric is worth telling someone.
type Rule struct {
	Name       string        // "pass-failures"
	Severity   Severity      // Critical pages; Warning goes on the morning list; Info is only logged
	FireAt     float64       // the problem starts when the value is at or above this
	ClearBelow float64       // and is over only when the value falls below this
	For        time.Duration // the value must stay at or above FireAt this long first
	Repeat     time.Duration // a problem still firing is sent again after this long; 0: never
}

// Sample is one value of a rule's metric for one station.
type Sample struct {
	At      time.Time
	Station string
	Value   float64
}

// Alert is one message the pager sends.
type Alert struct {
	At      time.Time
	Rule    string
	Station string
	Kind    string // "FIRING" or "RESOLVED"
	Value   float64
	To      string // the on-call name, "morning list" or "log"
}

// String writes the alert as the pager's log shows it: "20:21 FIRING pass-failures LAN 2".
func (a Alert) String() string {
	return fmt.Sprintf("%s %s %s %s %g", a.At.UTC().Format("15:04"), a.Kind, a.Rule, a.Station, a.Value)
}

// newRules returns the pager's rules from Monday 8 January 2029, 08:00.
func newRules() []Rule {
	return []Rule{ // TODO: the new rules; these are the old rules' settings
		{Name: "pass-failures", Severity: Critical, FireAt: 3, ClearBelow: 3},
		{Name: "failed-pass", Severity: Critical, FireAt: 1, ClearBelow: 1},
		{Name: "disk-full", Severity: Critical, FireAt: 0.98, ClearBelow: 0.98, Repeat: 10 * time.Minute},
	}
}

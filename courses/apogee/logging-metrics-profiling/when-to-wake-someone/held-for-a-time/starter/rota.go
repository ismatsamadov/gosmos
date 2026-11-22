package main

import "time"

// Shift is one week on call; it starts on a Monday at 08:00 UTC.
type Shift struct {
	From time.Time
	Name string // "Nkem Halloran"
}

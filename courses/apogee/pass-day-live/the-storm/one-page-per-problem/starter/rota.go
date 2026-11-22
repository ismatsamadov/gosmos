package main

import "time"

// Shift is one week on call; it starts on a Monday at 08:00 UTC.
type Shift struct {
	From time.Time
	Name string // "Nkem Halloran"
}

// onCall returns who carries the pager in July 2029, oldest first: a week each from Monday
// 08:00 UTC. Nkem swapped weeks with No. 039, so the rehearsal's author carries it for the
// investors' week.
func onCall() []Shift {
	monday := func(day int) time.Time { return time.Date(2029, time.July, day, 8, 0, 0, 0, time.UTC) }
	return []Shift{
		{monday(9), "Tomás Ferreira-Lund"},
		{monday(16), "Nkem Halloran"},
		{monday(23), "No. 039"},
		{monday(30), "Mirela Dunn"},
	}
}

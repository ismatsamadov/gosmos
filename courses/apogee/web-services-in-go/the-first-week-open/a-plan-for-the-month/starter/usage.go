package main

import "time"

// usage is one calendar month's count of bookings, as usage.json keeps it.
type usage struct {
	Month time.Time      `json:"month"` // the month's first instant, UTC
	Used  map[string]int `json:"used"`  // bookings made in it, by customer ID
}

// take counts one booking for customer at now against a plan of plan bookings a month. A new month starts a new count; a plan used up refuses the booking with a 403 that says when the count starts again.
func (u *usage) take(customer string, plan int, now time.Time) error {
	return nil // TODO: start a new month's count, refuse a booking past the plan, count the rest
}

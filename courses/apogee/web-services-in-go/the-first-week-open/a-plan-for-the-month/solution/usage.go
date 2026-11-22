package main

import (
	"fmt"
	"net/http"
	"time"
)

// usage is one calendar month's count of bookings, as usage.json keeps it.
type usage struct {
	Month time.Time      `json:"month"` // the month's first instant, UTC
	Used  map[string]int `json:"used"`  // bookings made in it, by customer ID
}

// take counts one booking for customer at now against a plan of plan bookings a month. A new month starts a new count; a plan used up refuses the booking with a 403 that says when the count starts again.
func (u *usage) take(customer string, plan int, now time.Time) error {
	now = now.UTC()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if !month.Equal(u.Month) {
		u.Month, u.Used = month, map[string]int{}
	}
	if u.Used[customer] >= plan {
		next := month.AddDate(0, 1, 0).Format(time.RFC3339)
		return &APIError{Status: http.StatusForbidden, Code: "quota_exceeded", Message: fmt.Sprintf("plan of %d bookings a month used up: the count starts again at %s", plan, next)}
	}
	u.Used[customer]++
	return nil
}

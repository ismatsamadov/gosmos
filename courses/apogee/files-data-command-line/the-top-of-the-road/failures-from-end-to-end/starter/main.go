package main

import (
	"fmt"
	"path/filepath"
)

// planDay reads the settings and the bookings in dir and plans the day.
func planDay(dir, format string) (Plan, error) {
	settings, _ := readSettings(filepath.Join(dir, "settings.json"))
	bookings, _ := readBookingFile(filepath.Join(dir, "bookings."+format), format)
	stations, _ := loadStations()
	// TODO: return each step's error at once, wrapped with what it was reading, and validate before planning.
	_ = fmt.Errorf
	plan, _ := newPlan(settings, stations, bookings)
	return plan, nil
}

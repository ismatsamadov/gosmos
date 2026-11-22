package main

import (
	"fmt"
	"path/filepath"
)

// planDay reads the settings and the bookings in dir and plans the day.
func planDay(dir, format string) (Plan, error) {
	settings, err := readSettings(filepath.Join(dir, "settings.json"))
	if err != nil {
		return Plan{}, fmt.Errorf("settings: %w", err)
	}
	bookings, err := readBookingFile(filepath.Join(dir, "bookings."+format), format)
	if err != nil {
		return Plan{}, fmt.Errorf("bookings: %w", err)
	}
	stations, err := loadStations()
	if err != nil {
		return Plan{}, err
	}
	if err := validateBookings(bookings, stations); err != nil {
		return Plan{}, fmt.Errorf("bookings: %w", err)
	}
	return newPlan(settings, stations, bookings)
}

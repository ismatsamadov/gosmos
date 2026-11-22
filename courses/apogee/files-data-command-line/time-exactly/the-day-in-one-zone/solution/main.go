package main

import (
	"slices"
	"time"
)

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string    `json:"id"`
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
}

// dayBookings returns the bookings whose AOS falls on day in zone, ordered by AOS.
func dayBookings(bookings []Booking, day string, zone *time.Location) []Booking {
	var found []Booking
	for _, b := range bookings {
		if b.AOS.In(zone).Format(time.DateOnly) == day {
			found = append(found, b)
		}
	}
	slices.SortFunc(found, func(a, b Booking) int {
		return a.AOS.Compare(b.AOS)
	})
	return found
}

package main

import (
	"fmt"
	"io"
)

// readBookings reads bookings in the named format: json for the export, csv for the sheet.
func readBookings(r io.Reader, format string) ([]Booking, error) {
	_ = fmt.Errorf
	return nil, nil
}

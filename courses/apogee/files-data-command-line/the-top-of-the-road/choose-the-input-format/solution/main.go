package main

import (
	"fmt"
	"io"
)

// readBookings reads bookings in the named format: json for the export, csv for the sheet.
func readBookings(r io.Reader, format string) ([]Booking, error) {
	switch format {
	case "json":
		return decodeBookings(r)
	case "csv":
		return readBookingSheet(r)
	default:
		return nil, fmt.Errorf("unsupported booking format %q", format)
	}
}

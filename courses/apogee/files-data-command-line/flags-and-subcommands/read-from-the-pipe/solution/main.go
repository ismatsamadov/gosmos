package main

import (
	"io"
	"os"
)

// loadBookings reads the bookings in the named file, or on stdin when no file is named.
func loadBookings(name string, stdin io.Reader) ([]Booking, error) {
	if name == "" {
		return readBookings(stdin)
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readBookings(f)
}

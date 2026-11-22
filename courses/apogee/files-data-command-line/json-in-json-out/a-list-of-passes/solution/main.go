package main

import "encoding/json"

// decodeBookings reads a day's bookings from a JSON array; a day without any gives an empty list, never nil.
func decodeBookings(data []byte) ([]Booking, error) {
	var bookings []Booking
	if err := json.Unmarshal(data, &bookings); err != nil {
		return nil, err
	}
	if bookings == nil {
		bookings = []Booking{}
	}
	return bookings, nil
}

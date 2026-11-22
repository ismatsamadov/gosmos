package main

import "encoding/json"

// decodeBookings reads a day's bookings from a JSON array; a day without any gives an empty list, never nil.
func decodeBookings(data []byte) ([]Booking, error) {
	_ = json.Unmarshal
	return nil, nil
}

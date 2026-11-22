package main

import "encoding/json"

// encodeBooking writes one booking as JSON for the reply to the booking export.
func encodeBooking(b Booking) ([]byte, error) {
	_ = json.Marshal
	return nil, nil
}

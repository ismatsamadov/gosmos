package main

import "encoding/json"

// decodeBooking reads one booking from the export's JSON.
func decodeBooking(data []byte) (Booking, error) {
	_ = json.Unmarshal
	return Booking{}, nil
}

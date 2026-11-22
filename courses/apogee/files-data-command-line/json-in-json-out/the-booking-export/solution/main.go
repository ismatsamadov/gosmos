package main

import "encoding/json"

// decodeBooking reads one booking from the export's JSON.
func decodeBooking(data []byte) (Booking, error) {
	var b Booking
	if err := json.Unmarshal(data, &b); err != nil {
		return Booking{}, err
	}
	return b, nil
}

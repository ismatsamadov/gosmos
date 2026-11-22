package main

import "unicode/utf8"

// checkShape returns a 422 *APIError for the first rule b's fields break, or nil when they break none.
func checkShape(b Booking) error {
	switch {
	case b.ID != "":
		return invalidBooking("id is chosen by the service")
	case b.SatID == "":
		return invalidBooking("sat_id is required")
	case b.Station == "":
		return invalidBooking("station is required")
	case b.AOS.IsZero() || b.LOS.IsZero():
		return invalidBooking("aos and los are required")
	case !b.LOS.After(b.AOS):
		return invalidBooking("los must be after aos")
	case utf8.RuneCountInString(b.Note) > 200:
		return invalidBooking("note is longer than 200 characters")
	}
	return nil
}

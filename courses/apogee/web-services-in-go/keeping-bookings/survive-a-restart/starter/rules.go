package main

import (
	"fmt"
	"time"
	"unicode/utf8"
)

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

// checkPass returns a 422 *APIError when b is not a pass in the catalog that is still open for booking at now, or nil when it is.
func checkPass(b Booking, catalog []PassView, now time.Time) error {
	if !knownSat(b.SatID) {
		return invalidBooking(fmt.Sprintf("unknown satellite %q", b.SatID))
	}
	if !knownStation(b.Station) {
		return invalidBooking(fmt.Sprintf("unknown station %q", b.Station))
	}
	for _, p := range catalog {
		if p.SatID == b.SatID && p.Station == b.Station && p.AOS.Equal(b.AOS) && p.LOS.Equal(b.LOS) {
			closes := p.AOS.Add(-10 * time.Minute)
			if !now.Before(closes) {
				return invalidBooking("bookings for this pass closed at " + closes.UTC().Format(time.RFC3339))
			}
			return nil
		}
	}
	return invalidBooking(fmt.Sprintf("no predicted pass of %s over %s at %s", b.SatID, b.Station, b.AOS.UTC().Format(time.RFC3339)))
}

package main

import (
	"fmt"
	"time"
)

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

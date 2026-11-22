package main

import "fmt"

// validateBookings checks every booking, returning the first problem it finds.
func validateBookings(bookings []Booking, stations map[string]Station) error {
	for i, b := range bookings {
		_, known := stations[b.Station]
		switch {
		case b.ID == "":
			return fmt.Errorf("booking %d: no id", i+1)
		case b.AOS.IsZero():
			return fmt.Errorf("booking %s: no AOS", b.ID)
		case !b.LOS.After(b.AOS):
			return fmt.Errorf("booking %s: LOS is not after AOS", b.ID)
		case !known:
			return fmt.Errorf("booking %s: unknown station %q", b.ID, b.Station)
		}
	}
	return nil
}

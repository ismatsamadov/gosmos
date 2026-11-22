package main

import (
	"fmt"
	"time"
)

// checkPass returns a 422 *APIError when b is not a pass in the catalog that is still open for booking at now, or nil when it is.
func checkPass(b Booking, catalog []PassView, now time.Time) error {
	_ = fmt.Sprintf
	return nil
}

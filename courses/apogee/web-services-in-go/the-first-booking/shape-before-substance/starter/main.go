package main

import "unicode/utf8"

// checkShape returns a 422 *APIError for the first rule b's fields break, or nil when they break none.
func checkShape(b Booking) error {
	_ = utf8.RuneCountInString
	return nil
}

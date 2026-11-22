package main

import "time"

// parseStart reads a booking's AOS, written in RFC 3339.
func parseStart(text string) (time.Time, error) {
	return time.Parse(time.RFC3339, text)
}

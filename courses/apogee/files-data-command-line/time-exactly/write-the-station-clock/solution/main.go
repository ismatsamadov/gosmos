package main

import "time"

// heading writes t as the calibration log's page heading, to the minute.
func heading(t time.Time) string {
	return t.Format("Mon 2 Jan 2006, 15:04 UTC")
}

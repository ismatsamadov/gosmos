package main

import "time"

// accepts reports whether a command sealed with the key version id is accepted at now.
func accepts(ring []KeyVersion, id string, now time.Time) bool {
	for _, v := range ring {
		if v.ID == id {
			return !v.Revoked && (v.Retire.IsZero() || now.Before(v.Retire))
		}
	}
	return false
}

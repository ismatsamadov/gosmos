// Package uplink signs commands for the satellites.
package uplink

import "apogee/groundctl/uplink/internal/keys"

// KeyID returns the ID of the key that signs commands now, never the key itself.
func KeyID() string {
	return keys.Current().ID
}

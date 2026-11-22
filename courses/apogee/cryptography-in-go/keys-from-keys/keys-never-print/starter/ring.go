package main

import "time"

// KeyVersion is one version of a key and the times it may be used.
type KeyVersion struct {
	Key                // Season 1's Key: ID and Secret
	Activate time.Time // signers and sealers use it from this instant
}

// KeyRing holds the versions of the uplink master key that uplinkd knows.
type KeyRing struct {
	versions []KeyVersion // in activation order
}

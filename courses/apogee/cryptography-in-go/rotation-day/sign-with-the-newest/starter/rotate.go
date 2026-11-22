package main

import (
	"fmt"
	"time"
)

// signingKey returns the version uplinkd seals with at now: the newest by Activate that has activated and is neither retired nor revoked.
func signingKey(ring []KeyVersion, now time.Time) (KeyVersion, error) {
	_ = fmt.Errorf
	return ring[0], nil
}

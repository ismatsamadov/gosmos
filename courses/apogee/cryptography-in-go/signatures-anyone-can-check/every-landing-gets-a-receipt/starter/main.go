package main

import (
	"fmt"
	"time"
)

// issue signs the receipt for landing l of booking b, whose file has the given digest, at the
// time now gives, and returns its envelope.
func issue(n *Notary, l Landing, b Booking, digest string, now func() time.Time) ([]byte, error) {
	_, _ = fmt.Errorf, time.Second
	return nil, nil
}

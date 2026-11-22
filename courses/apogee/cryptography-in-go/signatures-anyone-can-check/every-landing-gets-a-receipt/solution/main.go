package main

import (
	"fmt"
	"time"
)

// issue signs the receipt for landing l of booking b, whose file has the given digest, at the
// time now gives, and returns its envelope.
func issue(n *Notary, l Landing, b Booking, digest string, now func() time.Time) ([]byte, error) {
	if l.BookingID != b.ID {
		return nil, fmt.Errorf("landing for %s, not %s", l.BookingID, b.ID)
	}
	if _, _, err := parseDigest(digest); err != nil {
		return nil, fmt.Errorf("receipt for %s: %w", b.ID, err)
	}
	r := Receipt{
		Version:   1,
		KeyID:     n.KeyID,
		BookingID: b.ID,
		SatID:     b.SatID,
		Station:   b.Station,
		File:      l.File,
		Bytes:     l.Bytes,
		Digest:    digest,
		LandedAt:  l.At,
		SignedAt:  now().UTC().Truncate(time.Second),
	}
	return wrapReceipt(n, r)
}

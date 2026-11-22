package main

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// KeyVersion is one version of a key and the times it may be used.
type KeyVersion struct {
	Key                // Season 1's Key: ID and Secret
	Activate time.Time // signers and sealers use it from this instant
	Retire   time.Time // nothing it made is accepted from this instant; zero: not scheduled
	Revoked  bool      // refused at once, whatever the times; never undone
}

// KeyRing holds the versions of the uplink master key that uplinkd knows.
type KeyRing struct {
	versions []KeyVersion // in activation order
}

// Versions returns a copy of the ring's versions, in activation order.
func (r *KeyRing) Versions() []KeyVersion {
	return slices.Clone(r.versions)
}

// String lists the ring's key IDs in activation order, and no secret.
func (r *KeyRing) String() string {
	ids := make([]string, len(r.versions))
	for i, v := range r.versions {
		ids[i] = v.ID
	}
	return "key ring: " + strings.Join(ids, ", ")
}

// GoString describes the ring as String does, for %#v.
func (r *KeyRing) GoString() string {
	return r.String()
}

// SatKey derives satID's key for purpose from the version of the master key named keyID.
func (r *KeyRing) SatKey(keyID, purpose, satID string) ([]byte, error) {
	for _, v := range r.versions {
		if v.ID == keyID {
			return satKey(v.Key, purpose, satID)
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrUnknownKey, keyID)
}

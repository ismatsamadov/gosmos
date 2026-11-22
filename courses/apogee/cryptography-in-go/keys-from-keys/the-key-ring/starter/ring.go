package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// ErrUnknownKey means the ring holds no version with the key ID asked for.
var ErrUnknownKey = errors.New("unknown key")

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

// add puts v into the ring in activation order, refusing an ID the ring already holds.
func (r *KeyRing) add(v KeyVersion) error {
	for _, have := range r.versions {
		if have.ID == v.ID {
			return fmt.Errorf("duplicate key ID %q", v.ID)
		}
	}
	r.versions = append(r.versions, v)
	slices.SortStableFunc(r.versions, func(a, b KeyVersion) int { return a.Activate.Compare(b.Activate) })
	return nil
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

// parseVersion reads one line of uplink-keys.txt: <id> <activate> <retire or -> <64 hex secret> [revoked].
func parseVersion(line string) (KeyVersion, error) {
	f := strings.Split(line, " ")
	if len(f) != 4 && (len(f) != 5 || f[4] != "revoked") {
		return KeyVersion{}, errors.New("want <id> <activate> <retire or -> <secret> [revoked]")
	}
	activate, err := time.Parse(time.RFC3339, f[1])
	if err != nil {
		return KeyVersion{}, fmt.Errorf("activate time: %w", err)
	}
	var retire time.Time
	if f[2] != "-" {
		if retire, err = time.Parse(time.RFC3339, f[2]); err != nil {
			return KeyVersion{}, fmt.Errorf("retire time: %w", err)
		}
	}
	secret, err := parseSecret(f[3])
	if err != nil {
		return KeyVersion{}, err
	}
	return KeyVersion{Key: Key{ID: f[0], Secret: secret}, Activate: activate, Retire: retire, Revoked: len(f) == 5}, nil
}

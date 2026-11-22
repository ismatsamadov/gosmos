package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

// loadRing reads uplink-keys.txt: one version of the master key per line, among blank lines and # notes.
func loadRing(r io.Reader) (*KeyRing, error) {
	ring := &KeyRing{}
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		v, err := parseVersion(line)
		if err == nil {
			err = ring.add(v)
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", n, err)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return ring, nil
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

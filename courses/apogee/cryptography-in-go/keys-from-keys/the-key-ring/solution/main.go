package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

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

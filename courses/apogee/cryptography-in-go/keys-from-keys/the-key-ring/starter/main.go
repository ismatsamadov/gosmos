package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// loadRing reads uplink-keys.txt: one version of the master key per line, among blank lines and # notes.
func loadRing(r io.Reader) (*KeyRing, error) {
	_, _ = bufio.NewScanner, strings.HasPrefix
	return &KeyRing{}, nil
}

// SatKey derives satID's key for purpose from the version of the master key named keyID.
func (r *KeyRing) SatKey(keyID, purpose, satID string) ([]byte, error) {
	_ = fmt.Errorf
	return nil, nil
}

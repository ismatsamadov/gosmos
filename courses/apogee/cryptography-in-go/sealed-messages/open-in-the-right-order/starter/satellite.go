package main

import (
	"fmt"
	"time"
)

// SatKeys are the keys a satellite holds for one master version.
type SatKeys struct {
	Sign, Seal []byte
	Retire     time.Time // refused from this instant; zero: not scheduled
	Revoked    bool
}

// Satellite stands in for a satellite's flight software: it opens the SEAL lines meant for it.
type Satellite struct {
	ID   string
	Keys map[string]SatKeys // by key ID
	Last uint64             // the last sequence number it accepted
	Now  func() time.Time
}

// keyFor returns the keys of master version id, or why the satellite won't use them now.
func (s *Satellite) keyFor(id string) (SatKeys, error) {
	k, ok := s.Keys[id]
	switch {
	case !ok:
		return SatKeys{}, fmt.Errorf("%w: %q", ErrUnknownKey, id)
	case k.Revoked:
		return SatKeys{}, fmt.Errorf("%w: %q", ErrRevoked, id)
	case !k.Retire.IsZero() && !s.Now().Before(k.Retire):
		return SatKeys{}, fmt.Errorf("%w: %q", ErrRetired, id)
	}
	return k, nil
}

package main

import (
	"crypto/ecdh"
	"crypto/rand"
)

// Backhaul is one station's end of the backhaul: it makes the X25519 key for each session.
type Backhaul struct {
	static *ecdh.PrivateKey // made at the first session, then kept in the config
}

// NewSession returns the X25519 key for one new session.
func (b *Backhaul) NewSession() (*ecdh.PrivateKey, error) {
	if b.static == nil {
		key, err := ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		b.static = key
	}
	return b.static, nil
}

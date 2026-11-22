package main

import (
	"crypto/ecdh"
	"crypto/rand"
)

// Backhaul is one station's end of the backhaul: it makes the X25519 key for each session.
type Backhaul struct{}

// NewSession returns a new X25519 key for one new session; the Backhaul keeps nothing of it.
func (b *Backhaul) NewSession() (*ecdh.PrivateKey, error) {
	return ecdh.X25519().GenerateKey(rand.Reader)
}

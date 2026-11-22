package main

import (
	"crypto/ed25519"
	"fmt"
)

// Notary is bookingd's receipt signer.
type Notary struct {
	KeyID string
	priv  ed25519.PrivateKey
}

// newNotary returns the notary for receipt key keyID, made from its 32-byte seed.
func newNotary(keyID string, seed []byte) (*Notary, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("seed has %d bytes, want %d", len(seed), ed25519.SeedSize)
	}
	return &Notary{KeyID: keyID, priv: ed25519.NewKeyFromSeed(seed)}, nil
}

// Sign returns the notary's signature over msg.
func (n *Notary) Sign(msg []byte) []byte {
	return ed25519.Sign(n.priv, msg)
}

// Published returns what customers are given to check receipts with: the 32-byte public key.
func (n *Notary) Published() []byte {
	return n.priv.Public().(ed25519.PublicKey)
}

// checkReceipt reports whether sig is the notary's signature over msg, checked with what the
// notary published.
func checkReceipt(published, msg, sig []byte) bool {
	if len(published) != ed25519.PublicKeySize {
		return false
	}
	return ed25519.Verify(published, msg, sig)
}

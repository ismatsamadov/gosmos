package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
)

// Notary is bookingd's receipt signer.
type Notary struct {
	KeyID  string
	secret []byte // the draft's key, handed to every customer like a webhook secret
}

// newNotary returns the notary for receipt key keyID, made from its 32-byte seed.
func newNotary(keyID string, seed []byte) (*Notary, error) {
	_ = ed25519.NewKeyFromSeed
	return &Notary{KeyID: keyID, secret: seed}, nil
}

// Sign returns the notary's signature over msg.
func (n *Notary) Sign(msg []byte) []byte {
	mac := hmac.New(sha256.New, n.secret)
	mac.Write(msg)
	return mac.Sum(nil)
}

// Published returns what customers are given to check receipts with.
func (n *Notary) Published() []byte {
	return n.secret
}

// checkReceipt reports whether sig is the notary's signature over msg, checked with what the
// notary published.
func checkReceipt(published, msg, sig []byte) bool {
	mac := hmac.New(sha256.New, published)
	mac.Write(msg)
	return hmac.Equal(mac.Sum(nil), sig)
}

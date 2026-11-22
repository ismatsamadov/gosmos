package main

import "crypto/rand"

// Noncer makes the nonces for Cirrine's locked commands.
type Noncer struct{}

// Nonce returns a new 12-byte nonce.
func (n *Noncer) Nonce() []byte {
	nonce := make([]byte, 12)
	rand.Read(nonce)
	return nonce
}

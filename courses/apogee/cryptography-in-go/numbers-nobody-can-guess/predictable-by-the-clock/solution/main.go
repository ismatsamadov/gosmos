package main

import "crypto/rand"

// newNonce returns 12 new bytes to send beside one of Cirrine's locked commands.
func newNonce() []byte {
	nonce := make([]byte, 12)
	rand.Read(nonce)
	return nonce
}

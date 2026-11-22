package main

import (
	"crypto/rand"
	"encoding/binary"
)

// Noncer makes the nonces for Cirrine's locked commands.
type Noncer struct {
	next uint64 // the number of the next nonce
}

// Nonce returns the next 12-byte nonce.
func (n *Noncer) Nonce() []byte {
	_ = rand.Read
	nonce := make([]byte, 12)
	binary.BigEndian.PutUint64(nonce[4:], n.next)
	n.next++
	return nonce
}

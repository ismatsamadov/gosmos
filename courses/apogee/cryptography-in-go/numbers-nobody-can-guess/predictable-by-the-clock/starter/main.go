package main

import (
	"math/rand/v2"
	"time"
)

// newNonce returns 12 new bytes to send beside one of Cirrine's locked commands.
func newNonce() []byte {
	r := rand.New(rand.NewPCG(uint64(time.Now().Unix()), 0))
	nonce := make([]byte, 12)
	for i := range nonce {
		nonce[i] = byte(r.Uint32())
	}
	return nonce
}

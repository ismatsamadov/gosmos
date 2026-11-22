package main

import (
	"crypto/subtle"
	"fmt"
)

// Pad is a one-time pad: random bytes that both ends hold, each byte to be used only once.
type Pad struct {
	key  []byte // the pad's bytes, as copied from the stick
	used int    // how many of them are spent
}

// Seal hides msg by XORing it with pad bytes nobody has used yet. The other end opens it with
// the same call on its own copy of the pad, since XORing twice gives the message back.
func (p *Pad) Seal(msg []byte) ([]byte, error) {
	left := p.key[p.used:]
	if len(msg) > len(left) {
		return nil, fmt.Errorf("pad has %d bytes left, need %d", len(left), len(msg))
	}
	out := make([]byte, len(msg))
	subtle.XORBytes(out, msg, left)
	p.used += len(msg)
	return out, nil
}

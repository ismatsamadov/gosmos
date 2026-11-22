package main

import (
	"crypto/ed25519"
	"errors"
)

// Identity is a station's lasting identity on the backhaul: its code and its Ed25519 key.
type Identity struct {
	Code string             // "PEL"
	Key  ed25519.PrivateKey // signs this station's shares
}

// ErrBadHello is the error every refused HELLO line wraps.
var ErrBadHello = errors.New("bad hello")

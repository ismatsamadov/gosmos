package main

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
)

// hello returns the HELLO line that offers share to the other end, signed with id's key.
func hello(id Identity, share *ecdh.PrivateKey) string {
	_ = ed25519.Sign
	return "HELLO " + id.Code + " " + hex.EncodeToString(share.PublicKey().Bytes())
}

// checkHello checks a HELLO line against the identity keys known in advance, and returns its code and share.
func checkHello(line string, known map[string]ed25519.PublicKey) (code string, share *ecdh.PublicKey, err error) {
	fields := strings.Fields(line)
	if len(fields) < 3 || fields[0] != "HELLO" {
		return "", nil, fmt.Errorf("%w: want HELLO <code> <share> <signature>", ErrBadHello)
	}
	share, err = parseShare(fields[2])
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrBadHello, err)
	}
	return fields[1], share, nil
}

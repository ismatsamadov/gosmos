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
	shareHex := hex.EncodeToString(share.PublicKey().Bytes())
	sig := ed25519.Sign(id.Key, []byte("backhaul share|"+id.Code+"|"+shareHex))
	return "HELLO " + id.Code + " " + shareHex + " " + hex.EncodeToString(sig)
}

// checkHello checks a HELLO line against the identity keys known in advance, and returns its code and share.
func checkHello(line string, known map[string]ed25519.PublicKey) (code string, share *ecdh.PublicKey, err error) {
	fields := strings.Fields(line)
	if len(fields) != 4 || fields[0] != "HELLO" {
		return "", nil, fmt.Errorf("%w: want HELLO <code> <share> <signature>", ErrBadHello)
	}
	code = fields[1]
	pub, ok := known[code]
	if !ok {
		return "", nil, fmt.Errorf("%w: unknown station %q", ErrBadHello, code)
	}
	share, err = parseShare(fields[2])
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrBadHello, err)
	}
	sig, err := hex.DecodeString(fields[3])
	if err != nil || !ed25519.Verify(pub, []byte("backhaul share|"+code+"|"+fields[2]), sig) {
		return "", nil, fmt.Errorf("%w: signature does not verify", ErrBadHello)
	}
	return code, share, nil
}

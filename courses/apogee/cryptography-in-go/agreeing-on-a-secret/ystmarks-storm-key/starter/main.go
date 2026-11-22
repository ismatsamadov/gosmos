package main

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
)

// dial runs Port Ellery's side of the handshake with the station peer over conn, and returns the session key.
func dial(conn io.ReadWriter, me Identity, peer string, known map[string]ed25519.PublicKey) ([]byte, error) {
	_, _, _ = ecdh.X25519, rand.Reader, fmt.Errorf
	return nil, nil
}

// answer runs the answering station's side of the handshake over conn, and returns the dialer's code and the session key.
func answer(conn io.ReadWriter, me Identity, known map[string]ed25519.PublicKey) (peer string, key []byte, err error) {
	return "", nil, nil
}

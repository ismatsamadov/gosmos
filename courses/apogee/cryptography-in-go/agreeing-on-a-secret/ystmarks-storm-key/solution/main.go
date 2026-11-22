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
	mine, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	if err := writeLine(conn, hello(me, mine)); err != nil {
		return nil, err
	}
	code, theirs, err := readHello(conn, known)
	if err != nil {
		return nil, err
	}
	if code != peer {
		return nil, fmt.Errorf("%w: %s answered, want %s", ErrBadHello, code, peer)
	}
	shared, err := mine.ECDH(theirs)
	if err != nil {
		return nil, err
	}
	return sessionKey(shared, mine.PublicKey().Bytes(), theirs.Bytes())
}

// answer runs the answering station's side of the handshake over conn, and returns the dialer's code and the session key.
func answer(conn io.ReadWriter, me Identity, known map[string]ed25519.PublicKey) (peer string, key []byte, err error) {
	peer, theirs, err := readHello(conn, known)
	if err != nil {
		return "", nil, err
	}
	mine, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return "", nil, err
	}
	if err := writeLine(conn, hello(me, mine)); err != nil {
		return "", nil, err
	}
	shared, err := mine.ECDH(theirs)
	if err != nil {
		return "", nil, err
	}
	key, err = sessionKey(shared, theirs.Bytes(), mine.PublicKey().Bytes())
	return peer, key, err
}

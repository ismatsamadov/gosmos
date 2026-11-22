package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// publishKey writes pub as a PEM PUBLIC KEY block, as the customer guide prints it.
func publishKey(pub ed25519.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}

// readKey reads an Ed25519 public key back from one PEM PUBLIC KEY block.
func readKey(pemText []byte) (ed25519.PublicKey, error) {
	block, rest := pem.Decode(pemText)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("want one PUBLIC KEY block")
	}
	if len(bytes.TrimSpace(rest)) != 0 {
		return nil, errors.New("text after the PUBLIC KEY block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("not an Ed25519 key")
	}
	return pub, nil
}

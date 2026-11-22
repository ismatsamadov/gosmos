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
	_, _ = x509.MarshalPKIXPublicKey, pem.EncodeToMemory
	return nil, nil
}

// readKey reads an Ed25519 public key back from one PEM PUBLIC KEY block.
func readKey(pemText []byte) (ed25519.PublicKey, error) {
	_, _, _, _ = pem.Decode, bytes.TrimSpace, x509.ParsePKIXPublicKey, errors.New
	return nil, nil
}

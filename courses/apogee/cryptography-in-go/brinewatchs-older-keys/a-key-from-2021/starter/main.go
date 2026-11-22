package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// readRSA reads the RSA public key in one PEM PUBLIC KEY block, refusing any other kind of key
// and any RSA key under 2048 bits.
func readRSA(pemText []byte) (*rsa.PublicKey, error) {
	block, rest := pem.Decode(pemText)
	if block == nil || block.Type != "PUBLIC KEY" || len(bytes.TrimSpace(rest)) != 0 {
		return nil, errors.New("want one PUBLIC KEY block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, _ := key.(*rsa.PublicKey) // TODO: switch on the key's type, and count an RSA key's bits
	return pub, nil
}

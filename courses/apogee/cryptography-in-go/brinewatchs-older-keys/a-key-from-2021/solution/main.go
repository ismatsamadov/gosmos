package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
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
	switch k := key.(type) {
	case *rsa.PublicKey:
		if k.N.BitLen() < 2048 {
			return nil, fmt.Errorf("want at least 2048 bits, got %d", k.N.BitLen())
		}
		return k, nil
	case ed25519.PublicKey:
		return nil, errors.New("want an RSA key, got Ed25519")
	case *ecdh.PublicKey:
		return nil, errors.New("want an RSA key, got X25519")
	default:
		return nil, fmt.Errorf("want an RSA key, got %T", k)
	}
}

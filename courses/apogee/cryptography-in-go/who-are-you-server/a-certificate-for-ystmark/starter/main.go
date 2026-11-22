package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"time"
)

// issueServer issues host's server certificate for the public key pub, signed by the CA ca.
func issueServer(ca Issuer, host string, serial int64, pub ed25519.PublicKey, now time.Time) (*x509.Certificate, error) {
	_, _, _, _, _ = rand.Reader, pkix.Name{}, big.NewInt, x509.CreateCertificate, errors.New
	return nil, nil
}

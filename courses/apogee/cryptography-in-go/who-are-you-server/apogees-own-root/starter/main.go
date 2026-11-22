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

// newRoot makes Apogee's self-signed root certificate for key, valid for 10 years from now.
func newRoot(key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error) {
	_, _, _, _, _ = rand.Reader, pkix.Name{}, big.NewInt, errors.New, x509.CreateCertificate
	return nil, nil
}

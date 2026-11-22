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

// Issuer is a certificate authority: its certificate and the key it signs with.
type Issuer struct {
	Cert *x509.Certificate
	Key  ed25519.PrivateKey
}

// newRoot makes Apogee's self-signed root certificate for key, valid for 10 years from now.
func newRoot(key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error) {
	if len(key) != ed25519.PrivateKeySize {
		return nil, errors.New("root key must be an Ed25519 private key")
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Apogee Root CA 2028", Organization: []string{"Apogee"}},
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),
		IsCA:                  true,
		BasicConstraintsValid: true,
		MaxPathLen:            1,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(der)
}

// newIntermediate makes the stations CA's certificate for key, signed by root, valid for 2 years from now.
func newIntermediate(root Issuer, key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error) {
	if len(key) != ed25519.PrivateKeySize || len(root.Key) != ed25519.PrivateKeySize || root.Cert == nil {
		return nil, errors.New("stations CA needs its own key and the root's certificate and key")
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Apogee Stations CA 2028", Organization: []string{"Apogee"}},
		NotBefore:             now,
		NotAfter:              now.AddDate(2, 0, 0),
		IsCA:                  true,
		BasicConstraintsValid: true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, root.Cert, key.Public(), root.Key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(der)
}

package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"path/filepath"
	"time"
)

// issueFiles makes h a new key and its certificate, and returns h's two files: the certificate chain, then the key, as PEM.
func issueFiles(ca Issuer, h Holder, now time.Time) (certPEM, keyPEM []byte, err error) {
	_, _, _, _ = ed25519.GenerateKey, rand.Reader, x509.MarshalPKCS8PrivateKey, pem.EncodeToMemory
	return nil, nil, nil
}

// issueAll issues the certificate of every holder for now, and saves each, with its key, in dir.
func issueAll(ca Issuer, dir string, now time.Time) error {
	_, _ = fmt.Errorf, filepath.Join
	return nil
}

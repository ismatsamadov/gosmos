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
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	issue := issueServer
	if h.Client {
		issue = issueClient
	}
	cert, err := issue(ca, h.CN, h.Serial, pub, now)
	if err != nil {
		return nil, nil, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	certPEM = append(certPEM, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.Cert.Raw})...)
	return certPEM, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// issueAll issues the certificate of every holder for now, and saves each, with its key, in dir.
func issueAll(ca Issuer, dir string, now time.Time) error {
	certs := make([][]byte, len(holders))
	keys := make([][]byte, len(holders))
	for i, h := range holders {
		var err error
		if certs[i], keys[i], err = issueFiles(ca, h, now); err != nil {
			return fmt.Errorf("%s: %w", h.Name, err)
		}
	}
	for i, h := range holders {
		if err := saveFile(filepath.Join(dir, h.Name+".pem"), certs[i], 0o644); err != nil {
			return err
		}
		if err := saveFile(filepath.Join(dir, h.Name+".key"), keys[i], 0o600); err != nil {
			return err
		}
	}
	return nil
}

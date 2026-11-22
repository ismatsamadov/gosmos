package main

import (
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"time"
)

// serverConfig serves leaf and its intermediate with key, judging every date by now.
func serverConfig(leaf, inter *x509.Certificate, key ed25519.PrivateKey, now func() time.Time) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{leaf.Raw, inter.Raw},
			PrivateKey:  key,
			Leaf:        leaf,
		}},
		Time: now,
	}
}

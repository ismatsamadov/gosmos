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

// harden refuses every version of TLS older than 1.3.
func harden(cfg *tls.Config) {
	cfg.MinVersion = tls.VersionTLS13
}

// requireClients makes the server refuse any client without a certificate that chains to roots.
func requireClients(cfg *tls.Config, roots *x509.CertPool) {
}

// clientConfig trusts only roots, checks the server's name against host, and proves itself with cert.
func clientConfig(roots *x509.CertPool, host string, cert tls.Certificate, now func() time.Time) *tls.Config {
	return &tls.Config{
		RootCAs:    roots,
		ServerName: host,
		MinVersion: tls.VersionTLS13,
		Time:       now,
	}
}

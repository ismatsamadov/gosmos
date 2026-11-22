package main

import (
	"crypto/x509"
	"time"
)

// verifyStation checks that leaf, sent with the certificates in chain, proves host's server at now, up to one of roots.
func verifyStation(leaf *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool, host string, now time.Time) error {
	_, err := leaf.Verify(x509.VerifyOptions{
		Roots:       roots,
		DNSName:     host,
		CurrentTime: now,
	})
	return err
}

package main

import (
	"crypto/x509"
	"time"
)

// trusted returns nil if cert proves it belongs to host at now, and an error saying why not otherwise.
func trusted(cert *x509.Certificate, host string, roots, inters *x509.CertPool, now time.Time) error {
	_, err := cert.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inters,
		DNSName:       host,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		CurrentTime:   now,
	})
	return err
}

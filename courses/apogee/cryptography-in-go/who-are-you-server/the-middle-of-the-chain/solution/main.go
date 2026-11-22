package main

import (
	"crypto/x509"
	"time"
)

// verifyStation checks that leaf, sent with the certificates in chain, proves host's server at now, up to one of roots.
func verifyStation(leaf *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool, host string, now time.Time) error {
	inters := x509.NewCertPool()
	for _, cert := range chain {
		inters.AddCert(cert)
	}
	_, err := leaf.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inters,
		DNSName:       host,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		CurrentTime:   now,
	})
	return err
}

package main

import (
	"crypto/x509"
	"time"
)

// checkAt checks that leaf, sent with chain, proves host's server up to one of roots, at the time now gives.
func checkAt(leaf *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool, host string, now func() time.Time) error {
	inters := x509.NewCertPool()
	for _, cert := range chain {
		inters.AddCert(cert)
	}
	_, err := leaf.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inters,
		DNSName:       host,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		CurrentTime:   now(),
	})
	return err
}

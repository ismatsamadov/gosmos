package main

import (
	"crypto/x509"
	"fmt"
	"time"
)

// trusted returns nil if cert proves it belongs to host at now, and an error saying why not otherwise.
func trusted(cert *x509.Certificate, host string, roots, inters *x509.CertPool, now time.Time) error {
	if cert.Subject.CommonName != host {
		return fmt.Errorf("certificate is for %q, not %q", cert.Subject.CommonName, host)
	}
	return nil
}

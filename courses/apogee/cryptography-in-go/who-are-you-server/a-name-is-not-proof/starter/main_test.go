package main

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
	"time"
)

// loadCerts returns the certificates of the PEM file name, in order.
func loadCerts(t *testing.T, name string) []*x509.Certificate {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	var certs []*x509.Certificate
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			t.Fatalf("parsing %s: %v", name, err)
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		t.Fatalf("%s holds no certificate", name)
	}
	return certs
}

// poolOf returns a pool holding certs.
func poolOf(certs ...*x509.Certificate) *x509.CertPool {
	pool := x509.NewCertPool()
	for _, c := range certs {
		pool.AddCert(c)
	}
	return pool
}

// TestTrusted checks Ystmark's certificate and Mirela's, which carries the same name. roots.pem
// holds Apogee's root; yst.pem, Ystmark's certificate and then the stations CA's; impostor.pem,
// Mirela's certificate for yst.apogee.internal and then her own Apogee Root CA 2028.
func TestTrusted(t *testing.T) {
	roots := poolOf(loadCerts(t, "roots.pem")...)
	yst := loadCerts(t, "yst.pem")
	fake := loadCerts(t, "impostor.pem")
	issued := time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		chain []*x509.Certificate // the certificate, then whatever its sender adds
		host  string
		at    time.Time
		ok    bool
	}{
		{"Ystmark's certificate", yst, "yst.apogee.internal", issued, true},
		{"Mirela's certificate, from her own Apogee Root CA 2028", fake, "yst.apogee.internal", issued, false},
		{"Ystmark's certificate, for Port Ellery's name", yst, "pel.apogee.internal", issued, false},
		{"Ystmark's certificate, one second after it expires", yst, "yst.apogee.internal", yst[0].NotAfter.Add(time.Second), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := trusted(tt.chain[0], tt.host, roots, poolOf(tt.chain[1:]...), tt.at)
			if tt.ok && err != nil {
				t.Errorf("trusted(%s, %q) = %v, want nil", tt.name, tt.host, err)
			}
			if !tt.ok && err == nil {
				t.Errorf("trusted(%s, %q) = nil, want an error", tt.name, tt.host)
			}
		})
	}
}

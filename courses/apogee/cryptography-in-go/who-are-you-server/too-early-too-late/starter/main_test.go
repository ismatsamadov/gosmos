package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"testing"
	"time"
)

// The bible's days: the CAs are made on 18 September, the certificates issued from 22 September.
var (
	caDay    = time.Date(2028, 9, 18, 0, 0, 0, 0, time.UTC)
	issueDay = time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC)
)

// seedKey returns the test key whose seed is the SHA-256 of label.
func seedKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// TestEdges asks about Ystmark's certificate, valid from 2028-09-22T00:00:00Z to
// 2028-12-21T00:00:00Z, at each edge of that window, on a clock that says exactly what the test wants.
func TestEdges(t *testing.T) {
	must := func(cert *x509.Certificate, err error) *x509.Certificate {
		t.Helper()
		if err != nil {
			t.Fatalf("making a test certificate: %v", err)
		}
		return cert
	}
	rootKey, caKey := seedKey("ca root 2028"), seedKey("ca stations 2028")
	root := must(newRoot(rootKey, caDay))
	inter := must(newIntermediate(Issuer{root, rootKey}, caKey, caDay))
	ystPub := seedKey("cert yst.apogee.internal 1001").Public().(ed25519.PublicKey)
	yst := must(issueServer(Issuer{inter, caKey}, "yst.apogee.internal", 1001, ystPub, issueDay))
	roots := x509.NewCertPool()
	roots.AddCert(root)

	tests := []struct {
		name string
		at   time.Time
		ok   bool
	}{
		{"at NotBefore exactly", yst.NotBefore, true},
		{"one second before NotBefore", yst.NotBefore.Add(-time.Second), false},
		{"at NotAfter exactly", yst.NotAfter, true},
		{"one second after NotAfter", yst.NotAfter.Add(time.Second), false},
		{"at 00:20, on Lantern Atoll's clock running 40 minutes slow", yst.NotBefore.Add(20*time.Minute - 40*time.Minute), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := func() time.Time { return tt.at }
			err := checkAt(yst, []*x509.Certificate{inter}, roots, "yst.apogee.internal", clock)
			if tt.ok && err != nil {
				t.Errorf("checkAt(Ystmark's certificate, at %s) = %v, want nil", tt.at.Format(time.RFC3339), err)
			}
			if !tt.ok && err == nil {
				t.Errorf("checkAt(Ystmark's certificate, at %s) = nil, want an error", tt.at.Format(time.RFC3339))
			}
		})
	}
}

package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"testing"
	"time"
)

// madeOn is the day Mirela makes the root: Monday 18 September 2028, 00:00 UTC.
var madeOn = time.Date(2028, 9, 18, 0, 0, 0, 0, time.UTC)

// rootKey returns the test root key, whose seed is the SHA-256 of "ca root 2028".
func rootKey() ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("ca root 2028"))
	return ed25519.NewKeyFromSeed(seed[:])
}

// makeRoot runs newRoot with the test key, stopping the test unless it returns a certificate.
func makeRoot(t *testing.T) *x509.Certificate {
	t.Helper()
	cert, err := newRoot(rootKey(), madeOn)
	if err != nil {
		t.Fatalf("newRoot(root key, 2028-09-18) = %v, want a certificate", err)
	}
	if cert == nil {
		t.Fatalf("newRoot(root key, 2028-09-18) = nil, nil, want a certificate")
	}
	return cert
}

// TestFields reads the root's fields, as the table of Apogee's certificate authority gives them.
func TestFields(t *testing.T) {
	root := makeRoot(t)
	checks := []struct{ field, got, want string }{
		{"Subject.CommonName", root.Subject.CommonName, "Apogee Root CA 2028"},
		{"Subject.Organization", fmt.Sprint(root.Subject.Organization), "[Apogee]"},
		{"Issuer.CommonName", root.Issuer.CommonName, "Apogee Root CA 2028"},
		{"SerialNumber", fmt.Sprint(root.SerialNumber), "1"},
		{"IsCA", fmt.Sprint(root.IsCA), "true"},
		{"BasicConstraintsValid", fmt.Sprint(root.BasicConstraintsValid), "true"},
		{"MaxPathLen", fmt.Sprint(root.MaxPathLen), "1"},
		{"KeyUsage has KeyUsageCertSign", fmt.Sprint(root.KeyUsage&x509.KeyUsageCertSign != 0), "true"},
		{"KeyUsage has KeyUsageCRLSign", fmt.Sprint(root.KeyUsage&x509.KeyUsageCRLSign != 0), "true"},
		{"NotBefore", root.NotBefore.Format(time.RFC3339), "2028-09-18T00:00:00Z"},
		{"NotAfter", root.NotAfter.Format(time.RFC3339), "2038-09-18T00:00:00Z"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("newRoot(root key, 2028-09-18).%s = %s, want %s", c.field, c.got, c.want)
		}
	}
}

// TestSelfSigned checks that the root carries the key it was made with, signed by that key, and
// that it verifies as the only root in its own pool.
func TestSelfSigned(t *testing.T) {
	root := makeRoot(t)
	if pub, ok := root.PublicKey.(ed25519.PublicKey); !ok || !pub.Equal(rootKey().Public()) {
		t.Errorf("newRoot(root key, …).PublicKey = %v, want the root key's public half", root.PublicKey)
	}
	if err := root.CheckSignatureFrom(root); err != nil {
		t.Errorf("root.CheckSignatureFrom(root) = %v, want nil: a root signs itself", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(root)
	if _, err := root.Verify(x509.VerifyOptions{Roots: pool, CurrentTime: madeOn}); err != nil {
		t.Errorf("root.Verify(its own pool, 2028-09-18) = %v, want nil", err)
	}
}

// TestNoKey hands newRoot no key at all: an error, never a panic.
func TestNoKey(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("newRoot(nil, 2028-09-18) panicked: %v, want an error", r)
		}
	}()
	if root, err := newRoot(nil, madeOn); err == nil {
		got := "nil"
		if root != nil {
			got = "a certificate"
		}
		t.Errorf("newRoot(nil, 2028-09-18) = %s, nil, want an error", got)
	}
}

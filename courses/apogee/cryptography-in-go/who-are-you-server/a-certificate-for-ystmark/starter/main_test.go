package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"testing"
	"time"
)

// The bible's days: the root and the stations CA are made on 18 September, the stations'
// certificates issued from 22 September.
var (
	caDay    = time.Date(2028, 9, 18, 0, 0, 0, 0, time.UTC)
	issueDay = time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC)
)

// seedKey returns the test key whose seed is the SHA-256 of label.
func seedKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// publicOf returns the public half of the test key made from label.
func publicOf(label string) ed25519.PublicKey {
	return seedKey(label).Public().(ed25519.PublicKey)
}

// stationsCA makes Apogee's test root and stations CA with the given newRoot and newIntermediate.
func stationsCA(t *testing.T) (root *x509.Certificate, ca Issuer) {
	t.Helper()
	rootKey, caKey := seedKey("ca root 2028"), seedKey("ca stations 2028")
	root, err := newRoot(rootKey, caDay)
	if err != nil {
		t.Fatalf("newRoot = %v", err)
	}
	inter, err := newIntermediate(Issuer{Cert: root, Key: rootKey}, caKey, caDay)
	if err != nil {
		t.Fatalf("newIntermediate = %v", err)
	}
	return root, Issuer{Cert: inter, Key: caKey}
}

// issue runs issueServer for host with its test key, stopping the test unless it returns a certificate.
func issue(t *testing.T, ca Issuer, host string, serial int64) *x509.Certificate {
	t.Helper()
	label := fmt.Sprintf("cert %s %d", host, serial)
	cert, err := issueServer(ca, host, serial, publicOf(label), issueDay)
	if err != nil {
		t.Fatalf("issueServer(stations CA, %q, %d, …, 2028-09-22) = %v, want a certificate", host, serial, err)
	}
	if cert == nil {
		t.Fatalf("issueServer(stations CA, %q, %d, …, 2028-09-22) = nil, nil, want a certificate", host, serial)
	}
	return cert
}

// check verifies cert for host and usage at at, with Apogee's root and the stations CA.
func check(cert, root *x509.Certificate, ca Issuer, host string, usage x509.ExtKeyUsage, at time.Time) ([][]*x509.Certificate, error) {
	roots, inters := x509.NewCertPool(), x509.NewCertPool()
	roots.AddCert(root)
	inters.AddCert(ca.Cert)
	return cert.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inters,
		DNSName:       host,
		KeyUsages:     []x509.ExtKeyUsage{usage},
		CurrentTime:   at,
	})
}

// TestFields reads Ystmark's certificate, the first the stations CA issues.
func TestFields(t *testing.T) {
	_, ca := stationsCA(t)
	yst := issue(t, ca, "yst.apogee.internal", 1001)
	checks := []struct{ field, got, want string }{
		{"Subject.CommonName", yst.Subject.CommonName, "yst.apogee.internal"},
		{"Subject.Organization", fmt.Sprint(yst.Subject.Organization), "[Apogee]"},
		{"DNSNames", fmt.Sprint(yst.DNSNames), "[yst.apogee.internal]"},
		{"Issuer.CommonName", yst.Issuer.CommonName, "Apogee Stations CA 2028"},
		{"SerialNumber", fmt.Sprint(yst.SerialNumber), "1001"},
		{"NotBefore", yst.NotBefore.Format(time.RFC3339), "2028-09-22T00:00:00Z"},
		{"NotAfter", yst.NotAfter.Format(time.RFC3339), "2028-12-21T00:00:00Z"},
		{"IsCA", fmt.Sprint(yst.IsCA), "false"},
		{"KeyUsage is KeyUsageDigitalSignature", fmt.Sprint(yst.KeyUsage == x509.KeyUsageDigitalSignature), "true"},
		{"ExtKeyUsage", fmt.Sprint(yst.ExtKeyUsage), fmt.Sprint([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth})},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("issueServer(stations CA, \"yst.apogee.internal\", 1001, …).%s = %s, want %s", c.field, c.got, c.want)
		}
	}
	if pub, ok := yst.PublicKey.(ed25519.PublicKey); !ok || !pub.Equal(publicOf("cert yst.apogee.internal 1001")) {
		t.Errorf("issueServer(…).PublicKey = %v, want Ystmark's public key", yst.PublicKey)
	}
}

// TestChain follows Ystmark's certificate through the stations CA to Apogee's root.
func TestChain(t *testing.T) {
	root, ca := stationsCA(t)
	yst := issue(t, ca, "yst.apogee.internal", 1001)
	chains, err := check(yst, root, ca, "yst.apogee.internal", x509.ExtKeyUsageServerAuth, issueDay)
	if err != nil {
		t.Fatalf("Verify(Ystmark's certificate, server auth, 2028-09-22) = %v, want a chain to the root", err)
	}
	var names []string
	for _, c := range chains[0] {
		names = append(names, c.Subject.CommonName)
	}
	if got, want := fmt.Sprint(names), "[yst.apogee.internal Apogee Stations CA 2028 Apogee Root CA 2028]"; got != want {
		t.Errorf("Verify(Ystmark's certificate) chain = %s, want %s", got, want)
	}
}

// TestNotForClients asks Ystmark's certificate to prove a client: it proves servers only.
func TestNotForClients(t *testing.T) {
	root, ca := stationsCA(t)
	yst := issue(t, ca, "yst.apogee.internal", 1001)
	if _, err := check(yst, root, ca, "yst.apogee.internal", x509.ExtKeyUsageClientAuth, issueDay); err == nil {
		t.Errorf("Verify(Ystmark's certificate, client auth) = nil, want an error")
	}
}

// TestOtherHost issues Port Ellery's certificate: it carries PEL's name, and proves nothing for YST.
func TestOtherHost(t *testing.T) {
	root, ca := stationsCA(t)
	pel := issue(t, ca, "pel.apogee.internal", 1002)
	if _, err := check(pel, root, ca, "pel.apogee.internal", x509.ExtKeyUsageServerAuth, issueDay); err != nil {
		t.Errorf("Verify(Port Ellery's certificate, %q) = %v, want nil", "pel.apogee.internal", err)
	}
	if _, err := check(pel, root, ca, "yst.apogee.internal", x509.ExtKeyUsageServerAuth, issueDay); err == nil {
		t.Errorf("Verify(Port Ellery's certificate, %q) = nil, want an error", "yst.apogee.internal")
	}
}

// TestNoIssuer hands issueServer the stations CA without its key, then without its certificate:
// each is an error, never a panic.
func TestNoIssuer(t *testing.T) {
	_, ca := stationsCA(t)
	pub := publicOf("cert yst.apogee.internal 1001")
	tests := []struct {
		name string
		ca   Issuer
	}{
		{"the stations CA without its key", Issuer{Cert: ca.Cert}},
		{"the stations CA without its certificate", Issuer{Key: ca.Key}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("issueServer(%s, …) panicked: %v, want an error", tt.name, r)
				}
			}()
			if cert, err := issueServer(tt.ca, "yst.apogee.internal", 1001, pub, issueDay); err == nil {
				got := "nil"
				if cert != nil {
					got = "a certificate"
				}
				t.Errorf("issueServer(%s, …) = %s, nil, want an error", tt.name, got)
			}
		})
	}
}

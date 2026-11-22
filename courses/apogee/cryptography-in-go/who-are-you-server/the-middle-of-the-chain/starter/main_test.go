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

// publicOf returns the public half of the test key made from label.
func publicOf(label string) ed25519.PublicKey {
	return seedKey(label).Public().(ed25519.PublicKey)
}

// TestVerifyStation sends the certificates a station might present, made with ca.go from the
// test seeds. Mirela's impostor uses one key, "ca someone else", for its own Apogee Root CA 2028
// and for the Apogee Stations CA 2028 under it.
func TestVerifyStation(t *testing.T) {
	must := func(cert *x509.Certificate, err error) *x509.Certificate {
		t.Helper()
		if err != nil {
			t.Fatalf("making a test certificate: %v", err)
		}
		return cert
	}
	rootKey, caKey, ystKey := seedKey("ca root 2028"), seedKey("ca stations 2028"), seedKey("cert yst.apogee.internal 1001")
	root := must(newRoot(rootKey, caDay))
	inter := must(newIntermediate(Issuer{root, rootKey}, caKey, caDay))
	stations := Issuer{inter, caKey}
	yst := must(issueServer(stations, "yst.apogee.internal", 1001, ystKey.Public().(ed25519.PublicKey), issueDay))
	uplinkd := must(issueClient(stations, "uplinkd", 2001, publicOf("cert uplinkd 2001"), issueDay))
	signedByYST := must(issueServer(Issuer{yst, ystKey}, "pel.apogee.internal", 1002, publicOf("cert impostor"), issueDay))
	fakeKey := seedKey("ca someone else")
	fakeRoot := must(newRoot(fakeKey, caDay))
	fakeInter := must(newIntermediate(Issuer{fakeRoot, fakeKey}, fakeKey, caDay))
	fake := must(issueServer(Issuer{fakeInter, fakeKey}, "yst.apogee.internal", 1001, publicOf("cert impostor"), issueDay))

	roots := x509.NewCertPool()
	roots.AddCert(root)
	tests := []struct {
		name  string
		leaf  *x509.Certificate
		chain []*x509.Certificate
		host  string
		ok    bool
	}{
		{"Ystmark's certificate with its chain", yst, []*x509.Certificate{inter}, "yst.apogee.internal", true},
		{"Ystmark's certificate without its chain", yst, nil, "yst.apogee.internal", false},
		{"a certificate for PEL that Ystmark's key signed", signedByYST, []*x509.Certificate{yst, inter}, "pel.apogee.internal", false},
		{"Mirela's certificate, under an intermediate from her own root", fake, []*x509.Certificate{fakeInter}, "yst.apogee.internal", false},
		{"uplinkd's client certificate, as a station", uplinkd, []*x509.Certificate{inter}, "uplinkd.apogee.internal", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyStation(tt.leaf, tt.chain, roots, tt.host, issueDay)
			if tt.ok && err != nil {
				t.Errorf("verifyStation(%s, %q) = %v, want nil", tt.name, tt.host, err)
			}
			if !tt.ok && err == nil {
				t.Errorf("verifyStation(%s, %q) = nil, want an error", tt.name, tt.host)
			}
		})
	}
}

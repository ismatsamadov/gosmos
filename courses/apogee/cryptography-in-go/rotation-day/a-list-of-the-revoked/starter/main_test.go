package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"math/big"
	"testing"
	"time"
)

// reported is 10:12 UTC on Rotation Day, when the Cape Wren contractor called.
var reported = time.Date(2028, 10, 17, 10, 12, 0, 0, time.UTC)

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// pki is Apogee's stations CA and the impostor's root, which also calls itself Apogee Root CA 2028.
type pki struct {
	inter    Issuer
	impostor Issuer
}

// newPKI builds both from the season's labels.
func newPKI(t *testing.T) pki {
	t.Helper()
	caDay := time.Date(2028, 9, 18, 0, 0, 0, 0, time.UTC)
	rootKey, interKey, impKey := testKey("ca root 2028"), testKey("ca stations 2028"), testKey("ca someone else")
	root, err := newRoot(rootKey, caDay)
	if err != nil {
		t.Fatalf("newRoot: %v", err)
	}
	inter, err := newIntermediate(Issuer{Cert: root, Key: rootKey}, interKey, caDay)
	if err != nil {
		t.Fatalf("newIntermediate: %v", err)
	}
	imp, err := newRoot(impKey, caDay)
	if err != nil {
		t.Fatalf("newRoot: %v", err)
	}
	return pki{inter: Issuer{Cert: inter, Key: interKey}, impostor: Issuer{Cert: imp, Key: impKey}}
}

// capeWren issues Cape Wren's certificate serial, as the stations CA did on the day given.
func capeWren(t *testing.T, p pki, serial int64, day time.Time) *x509.Certificate {
	t.Helper()
	key := testKey(fmt.Sprintf("cert cwr.apogee.internal %d", serial))
	leaf, err := issueServer(p.inter, "cwr.apogee.internal", serial, key.Public().(ed25519.PublicKey), day)
	if err != nil {
		t.Fatalf("issueServer: %v", err)
	}
	return leaf
}

// listOf is a revocation list the test signs itself with ca at the given time, naming serials.
func listOf(t *testing.T, ca Issuer, at time.Time, serials ...int64) []byte {
	t.Helper()
	var entries []x509.RevocationListEntry
	for _, s := range serials {
		entries = append(entries, x509.RevocationListEntry{SerialNumber: big.NewInt(s), RevocationTime: at})
	}
	der, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number: big.NewInt(1), ThisUpdate: at, NextUpdate: at.AddDate(0, 0, 7), RevokedCertificateEntries: entries,
	}, ca.Cert, ca.Key)
	if err != nil {
		t.Fatalf("x509.CreateRevocationList: %v", err)
	}
	return der
}

func TestRevoke(t *testing.T) {
	p := newPKI(t)
	entries := []x509.RevocationListEntry{{SerialNumber: big.NewInt(1005), RevocationTime: reported}}
	der, err := revoke(p.inter, 1, entries, reported)
	if err != nil || der == nil {
		t.Fatalf("revoke(stations CA, 1, [1005], %s) = %d bytes, %v, want a signed list", reported.Format(time.RFC3339), len(der), err)
	}
	list, err := x509.ParseRevocationList(der)
	if err != nil {
		t.Fatalf("x509.ParseRevocationList(revoke(...)) = %v, want a list", err)
	}
	if err := list.CheckSignatureFrom(p.inter.Cert); err != nil {
		t.Errorf("the list's signature = %v, want the stations CA's", err)
	}
	if list.Number == nil || list.Number.Int64() != 1 {
		t.Errorf("list.Number = %v, want 1", list.Number)
	}
	if !list.ThisUpdate.Equal(reported) {
		t.Errorf("list.ThisUpdate = %s, want %s", list.ThisUpdate.Format(time.RFC3339), reported.Format(time.RFC3339))
	}
	if want := reported.AddDate(0, 0, 7); !list.NextUpdate.Equal(want) {
		t.Errorf("list.NextUpdate = %s, want %s", list.NextUpdate.Format(time.RFC3339), want.Format(time.RFC3339))
	}
	if n := len(list.RevokedCertificateEntries); n != 1 || list.RevokedCertificateEntries[0].SerialNumber.Int64() != 1005 {
		t.Errorf("list entries = %d, want 1: certificate 1005", n)
	}
}

func TestIsRevoked(t *testing.T) {
	p := newPKI(t)
	old := capeWren(t, p, 1005, time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC))
	current := capeWren(t, p, 1105, time.Date(2028, 10, 9, 0, 0, 0, 0, time.UTC))
	list1 := listOf(t, p.inter, reported, 1005)
	forged := listOf(t, p.impostor, reported, 1005)
	next := reported.AddDate(0, 0, 7)
	tests := []struct {
		name    string
		list    []byte
		cert    *x509.Certificate
		now     time.Time
		want    bool
		wantErr bool
	}{
		{"1005 on list 1", list1, old, reported, true, false},
		{"1105, Cape Wren's new certificate", list1, current, reported, false, false},
		{"1005 on a list the impostor signed", forged, old, reported, false, true},
		{"1105 on a list the impostor signed", forged, current, reported, false, true},
		{"1005 at the list's NextUpdate exactly", list1, old, next, true, false},
		{"1105 a second after NextUpdate", list1, current, next.Add(time.Second), false, true},
		{"1105 a second before ThisUpdate", list1, current, reported.Add(-time.Second), false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isRevoked(tt.list, p.inter.Cert, tt.cert, tt.now)
			switch {
			case tt.wantErr && err == nil:
				t.Errorf("isRevoked(list, stations CA, %v, %s) = %v, nil, want an error", tt.cert.SerialNumber, tt.now.Format(time.RFC3339), got)
			case !tt.wantErr && (err != nil || got != tt.want):
				t.Errorf("isRevoked(list, stations CA, %v, %s) = %v, %v, want %v, nil", tt.cert.SerialNumber, tt.now.Format(time.RFC3339), got, err, tt.want)
			}
		})
	}
}

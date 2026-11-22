package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
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

// testCA makes Apogee's test root and stations CA with ca.go, from the test seeds.
func testCA(t *testing.T) (root *x509.Certificate, ca Issuer) {
	t.Helper()
	rootKey, caKey := seedKey("ca root 2028"), seedKey("ca stations 2028")
	root, err := newRoot(rootKey, caDay)
	if err != nil {
		t.Fatalf("newRoot = %v", err)
	}
	inter, err := newIntermediate(Issuer{root, rootKey}, caKey, caDay)
	if err != nil {
		t.Fatalf("newIntermediate = %v", err)
	}
	return root, Issuer{inter, caKey}
}

// verify checks leaf for name and usage on 22 September, up to root through the stations CA.
func verify(leaf, root, inter *x509.Certificate, name string, usage x509.ExtKeyUsage) error {
	roots, inters := x509.NewCertPool(), x509.NewCertPool()
	roots.AddCert(root)
	inters.AddCert(inter)
	_, err := leaf.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inters,
		DNSName:       name,
		KeyUsages:     []x509.ExtKeyUsage{usage},
		CurrentTime:   issueDay,
	})
	return err
}

// list returns the names of the files in dir, joined by spaces.
func list(t *testing.T, dir string) (int, string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return len(names), strings.Join(names, " ")
}

const (
	server = x509.ExtKeyUsageServerAuth
	client = x509.ExtKeyUsageClientAuth
)

// TestIssueAll issues every certificate into a temporary folder, then loads and checks each pair.
func TestIssueAll(t *testing.T) {
	root, ca := testCA(t)
	dir := t.TempDir()
	if err := issueAll(ca, dir, issueDay); err != nil {
		t.Fatalf("issueAll(stations CA, dir, 2028-09-22) = %v, want nil", err)
	}
	if n, names := list(t, dir); n != 16 {
		t.Fatalf("issueAll(stations CA, dir, 2028-09-22) left %d files (%s), want 16: a .pem and a .key for each of 8 holders", n, names)
	}
	tests := []struct {
		name   string // the files' name
		dns    string
		serial int64
		usage  x509.ExtKeyUsage
		wrong  x509.ExtKeyUsage // the use it must not prove
		other  string           // a name it must not prove
	}{
		{"yst", "yst.apogee.internal", 1001, server, client, "pel.apogee.internal"},
		{"pel", "pel.apogee.internal", 1002, server, client, "yst.apogee.internal"},
		{"sqt", "sqt.apogee.internal", 1003, server, client, "yst.apogee.internal"},
		{"lan", "lan.apogee.internal", 1004, server, client, "yst.apogee.internal"},
		{"cwr", "cwr.apogee.internal", 1005, server, client, "yst.apogee.internal"},
		{"uplinkd", "uplinkd.apogee.internal", 2001, client, server, "passplan.apogee.internal"},
		{"passplan", "passplan.apogee.internal", 2002, client, server, "uplinkd.apogee.internal"},
		{"dishd", "dishd.apogee.internal", 2003, client, server, "uplinkd.apogee.internal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			certPEM, err := os.ReadFile(filepath.Join(dir, tt.name+".pem"))
			if err != nil {
				t.Fatalf("reading %s.pem: %v", tt.name, err)
			}
			keyPEM, err := os.ReadFile(filepath.Join(dir, tt.name+".key"))
			if err != nil {
				t.Fatalf("reading %s.key: %v", tt.name, err)
			}
			if block, _ := pem.Decode(keyPEM); block == nil || block.Type != "PRIVATE KEY" {
				t.Errorf("pem.Decode(%s.key) = no PRIVATE KEY block, want the key as one PKCS #8 PRIVATE KEY block", tt.name)
			}
			pair, err := tls.X509KeyPair(certPEM, keyPEM)
			if err != nil {
				t.Fatalf("tls.X509KeyPair(%s.pem, %s.key) = %v, want a pair", tt.name, tt.name, err)
			}
			if len(pair.Certificate) != 2 || !bytes.Equal(pair.Certificate[1], ca.Cert.Raw) {
				t.Fatalf("%s.pem holds %d certificates, want 2: its own, then the stations CA's", tt.name, len(pair.Certificate))
			}
			leaf, err := x509.ParseCertificate(pair.Certificate[0])
			if err != nil {
				t.Fatalf("parsing %s.pem: %v", tt.name, err)
			}
			if leaf.SerialNumber.Int64() != tt.serial {
				t.Errorf("serial of %s.pem = %v, want %d", tt.name, leaf.SerialNumber, tt.serial)
			}
			if err := verify(leaf, root, ca.Cert, tt.dns, tt.usage); err != nil {
				t.Errorf("Verify(%s.pem, %q, its own use) = %v, want nil", tt.name, tt.dns, err)
			}
			if err := verify(leaf, root, ca.Cert, tt.dns, tt.wrong); err == nil {
				t.Errorf("Verify(%s.pem, %q, the other use) = nil, want an error", tt.name, tt.dns)
			}
			if err := verify(leaf, root, ca.Cert, tt.other, tt.usage); err == nil {
				t.Errorf("Verify(%s.pem, %q) = nil, want an error", tt.name, tt.other)
			}
			if info, err := os.Stat(filepath.Join(dir, tt.name+".key")); err == nil && info.Mode().Perm() != 0o600 {
				t.Errorf("mode of %s.key = %v, want -rw-------", tt.name, info.Mode().Perm())
			}
		})
	}
}

// TestNoKeyNoFiles hands issueAll the stations CA's certificate without its key.
func TestNoKeyNoFiles(t *testing.T) {
	_, ca := testCA(t)
	dir := t.TempDir()
	err := issueAll(Issuer{Cert: ca.Cert}, dir, issueDay)
	if err == nil {
		t.Errorf("issueAll(a CA without its key, dir, 2028-09-22) = nil, want an error")
	}
	if n, names := list(t, dir); n != 0 {
		t.Errorf("issueAll(a CA without its key) left %d files (%s), want none", n, names)
	}
}

// TestNewKeys checks that every holder gets a key of its own, and a new one on every run: keys
// come from crypto/rand.
func TestNewKeys(t *testing.T) {
	_, ca := testCA(t)
	first, second := t.TempDir(), t.TempDir()
	if err := issueAll(ca, first, issueDay); err != nil {
		t.Fatalf("issueAll(stations CA, first dir, …) = %v, want nil", err)
	}
	owner := make(map[string]string) // a key file's text → the holder it belongs to
	for _, h := range holders {
		key, err := os.ReadFile(filepath.Join(first, h.Name+".key"))
		if err != nil {
			t.Fatalf("reading %s.key: %v", h.Name, err)
		}
		if other, ok := owner[string(key)]; ok {
			t.Errorf("%s.key = %s.key, want a key of its own for every holder", h.Name, other)
			continue
		}
		owner[string(key)] = h.Name
	}
	if err := issueAll(ca, second, issueDay); err != nil {
		t.Fatalf("issueAll(stations CA, second dir, …) = %v, want nil", err)
	}
	a, errA := os.ReadFile(filepath.Join(first, "yst.key"))
	b, errB := os.ReadFile(filepath.Join(second, "yst.key"))
	if errA != nil || errB != nil {
		t.Fatalf("issueAll(stations CA, …) twice left no yst.key, want one from each run")
	}
	if bytes.Equal(a, b) {
		t.Errorf("yst.key = the same key on two runs, want a new key every time")
	}
}

// TestFailureLeavesNothing makes the fourth certificate impossible to issue, with a serial below
// zero, which x509.CreateCertificate refuses: the three issued before it must not be saved either.
func TestFailureLeavesNothing(t *testing.T) {
	_, ca := testCA(t)
	all := holders
	defer func() { holders = all }()
	holders = []Holder{all[0], all[1], all[2], {Name: "bad", CN: "bad.apogee.internal", Serial: -1}}
	dir := t.TempDir()
	if err := issueAll(ca, dir, issueDay); err == nil {
		t.Errorf("issueAll(stations CA, dir, …) with serial -1 fourth = nil, want an error")
	}
	if n, names := list(t, dir); n != 0 {
		t.Errorf("issueAll(…) with serial -1 fourth left %d files (%s), want none: nothing is saved until every holder is issued", n, names)
	}
}

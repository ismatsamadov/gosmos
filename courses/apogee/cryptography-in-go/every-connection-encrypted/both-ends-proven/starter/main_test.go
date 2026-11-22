package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// wednesday is the tests' clock: Wednesday 27 September 2028, 09:00 UTC.
var wednesday = time.Date(2028, 9, 27, 9, 0, 0, 0, time.UTC)

// now returns the tests' clock.
func now() time.Time { return wednesday }

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// pki is Apogee's test CA and the impostor's root of the same name.
type pki struct {
	roots    *x509.CertPool // Apogee Root CA 2028, alone
	inter    Issuer         // Apogee Stations CA 2028
	impostor Issuer         // another root named Apogee Root CA 2028, with the key "ca someone else"
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
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return pki{roots: roots, inter: Issuer{Cert: inter, Key: interKey}, impostor: Issuer{Cert: imp, Key: impKey}}
}

// issue makes the certificate serial for name from ca on 22 September, its key made from label, a
// client's or a server's, and returns it with its chain.
func issue(t *testing.T, ca Issuer, name string, serial int64, label string, client bool) tls.Certificate {
	t.Helper()
	key := testKey(label)
	mint := issueServer
	if client {
		mint = issueClient
	}
	leaf, err := mint(ca, name, serial, key.Public().(ed25519.PublicKey), time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("issuing %s %d: %v", name, serial, err)
	}
	return tls.Certificate{Certificate: [][]byte{leaf.Raw, ca.Cert.Raw}, PrivateKey: key, Leaf: leaf}
}

// ystmark starts Ystmark's uplinkd, hardened, asking clients for certificates with requireClients.
func ystmark(t *testing.T, p pki) *httptest.Server {
	t.Helper()
	yst := issue(t, p.inter, "yst.apogee.internal", 1001, "cert yst.apogee.internal 1001", false)
	cfg := serverConfig(yst.Leaf, p.inter.Cert, yst.PrivateKey.(ed25519.PrivateKey), now)
	harden(cfg)
	requireClients(cfg, p.roots)
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok YST\n")
	}))
	srv.TLS = cfg
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv
}

// status asks the server at url for its status with cfg, and returns the answer's status code.
func status(cfg *tls.Config, url string) (int, error) {
	c := &http.Client{Timeout: time.Second, Transport: &http.Transport{TLSClientConfig: cfg}}
	resp, err := c.Get(url + "/status")
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return resp.StatusCode, nil
}

func TestUplinkd(t *testing.T) {
	p := newPKI(t)
	srv := ystmark(t, p)
	up := issue(t, p.inter, "uplinkd", 2001, "cert uplinkd 2001", true)
	if code, err := status(clientConfig(p.roots, "yst.apogee.internal", up, now), srv.URL); code != http.StatusOK {
		t.Errorf("GET /status as uplinkd = %d, %v, want 200, nil", code, err)
	}
}

// TestRefused plays three clients Ystmark must turn away, each with a config of the test's own.
func TestRefused(t *testing.T) {
	p := newPKI(t)
	srv := ystmark(t, p)
	tests := []struct {
		name  string
		certs []tls.Certificate
	}{
		{"no certificate", nil},
		{"uplinkd's name from the impostor's root", []tls.Certificate{issue(t, p.impostor, "uplinkd", 2001, "cert impostor", true)}},
		{"Ystmark's own server certificate", []tls.Certificate{issue(t, p.inter, "yst.apogee.internal", 1001, "cert yst.apogee.internal 1001", false)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &tls.Config{RootCAs: p.roots, ServerName: "yst.apogee.internal", Certificates: tt.certs, MinVersion: tls.VersionTLS13, Time: now}
			if code, err := status(cfg, srv.URL); err == nil {
				t.Errorf("GET /status with %s = %d, want a refused handshake", tt.name, code)
			}
		})
	}
}

func TestClientConfig(t *testing.T) {
	p := newPKI(t)
	up := issue(t, p.inter, "uplinkd", 2001, "cert uplinkd 2001", true)
	cfg := clientConfig(p.roots, "yst.apogee.internal", up, now)
	var names []string
	for _, c := range cfg.Certificates {
		if c.Leaf != nil {
			names = append(names, fmt.Sprintf("%s %v", c.Leaf.Subject.CommonName, c.Leaf.SerialNumber))
		}
	}
	if len(cfg.Certificates) != 1 || cfg.Certificates[0].Leaf != up.Leaf {
		t.Errorf("clientConfig(roots, %q, uplinkd's, now).Certificates = %q, want [uplinkd 2001]", "yst.apogee.internal", names)
	}
}

package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// tuesday is the tests' clock: Tuesday 26 September 2028, 09:00 UTC.
var tuesday = time.Date(2028, 9, 26, 9, 0, 0, 0, time.UTC)

// now returns the tests' clock.
func now() time.Time { return tuesday }

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// ystmark starts Ystmark's uplinkd with serverConfig, hardened, and returns it with a pool that
// holds Apogee's root alone.
func ystmark(t *testing.T) (*httptest.Server, *x509.CertPool) {
	t.Helper()
	caDay := time.Date(2028, 9, 18, 0, 0, 0, 0, time.UTC)
	rootKey, interKey, ystKey := testKey("ca root 2028"), testKey("ca stations 2028"), testKey("cert yst.apogee.internal 1001")
	root, err := newRoot(rootKey, caDay)
	if err != nil {
		t.Fatalf("newRoot: %v", err)
	}
	inter, err := newIntermediate(Issuer{Cert: root, Key: rootKey}, interKey, caDay)
	if err != nil {
		t.Fatalf("newIntermediate: %v", err)
	}
	leaf, err := issueServer(Issuer{Cert: inter, Key: interKey}, "yst.apogee.internal", 1001,
		ystKey.Public().(ed25519.PublicKey), time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("issueServer: %v", err)
	}
	cfg := serverConfig(leaf, inter, ystKey, now)
	harden(cfg)
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok YST\n")
	}))
	srv.TLS = cfg
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	t.Cleanup(srv.Close)
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return srv, roots
}

// client returns a client that trusts roots alone and asks for yst.apogee.internal, offering TLS
// versions up to max (0: up to the newest).
func client(roots *x509.CertPool, max uint16) *http.Client {
	return &http.Client{Timeout: time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{
		RootCAs:    roots,
		ServerName: "yst.apogee.internal",
		MaxVersion: max,
		Time:       now,
	}}}
}

// TestOldTool plays the old tool at Salar Quieto, which speaks TLS 1.2 and nothing newer.
func TestOldTool(t *testing.T) {
	srv, roots := ystmark(t)
	resp, err := client(roots, tls.VersionTLS12).Get(srv.URL + "/status")
	if err == nil {
		resp.Body.Close()
		t.Errorf("GET /status from a client limited to TLS 1.2 = %s over %s, want a refused handshake",
			resp.Status, tls.VersionName(resp.TLS.Version))
	}
}

func TestOrdinaryClient(t *testing.T) {
	srv, roots := ystmark(t)
	resp, err := client(roots, 0).Get(srv.URL + "/status")
	if err != nil {
		t.Fatalf("GET /status = %v, want 200", err)
	}
	resp.Body.Close()
	if v := resp.TLS.Version; v != tls.VersionTLS13 {
		t.Errorf("GET /status: version = %s, want TLS 1.3", tls.VersionName(v))
	}
	if c := resp.TLS.CurveID; c != tls.X25519MLKEM768 {
		t.Errorf("GET /status: key exchange = %v, want X25519MLKEM768", c)
	}
}

func TestKeepsTheRest(t *testing.T) {
	cfg := &tls.Config{
		Certificates: make([]tls.Certificate, 1),
		ServerName:   "yst.apogee.internal",
		NextProtos:   []string{"http/1.1"},
		Time:         now,
	}
	harden(cfg)
	if cfg.MinVersion != tls.VersionTLS13 {
		t.Errorf("after harden(cfg): cfg.MinVersion = %s, want TLS 1.3", tls.VersionName(cfg.MinVersion))
	}
	if len(cfg.Certificates) != 1 || cfg.ServerName != "yst.apogee.internal" || len(cfg.NextProtos) != 1 || cfg.Time == nil {
		t.Errorf("after harden(cfg): %d certificates, ServerName %q, NextProtos %q, Time set %v, want them as they were: 1, %q, [http/1.1], true",
			len(cfg.Certificates), cfg.ServerName, cfg.NextProtos, cfg.Time != nil, "yst.apogee.internal")
	}
}

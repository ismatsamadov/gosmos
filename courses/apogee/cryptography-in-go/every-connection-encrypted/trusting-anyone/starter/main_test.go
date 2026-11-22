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
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// monday is the tests' clock: Monday 25 September 2028, 09:00 UTC.
var monday = time.Date(2028, 9, 25, 9, 0, 0, 0, time.UTC)

// now returns the tests' clock.
func now() time.Time { return monday }

// sealedPing is Cirrine-1's sealed PING from the season's test vectors: today's command for Ystmark.
const sealedPing = "SEAL 2 cirrine-1 4471 k-2028-07 a1b2c3d4e5f60718293a4b5c c1f2a9798d38081d764bbcfb5e35e2598d2584a4"

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// serve starts an HTTPS server for h on 127.0.0.1 with cfg, or with httptest's own certificate when
// cfg is nil. Refused handshakes print nothing.
func serve(t *testing.T, h http.Handler, cfg *tls.Config) *httptest.Server {
	t.Helper()
	srv := httptest.NewUnstartedServer(h)
	srv.TLS = cfg
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv
}

// ystmark starts Ystmark's uplinkd with the certificate Apogee's CA issued it on 22 September, and
// returns it with a pool holding Apogee's root alone.
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
	srv := serve(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok YST\n")
	}), &tls.Config{
		Certificates: []tls.Certificate{{Certificate: [][]byte{leaf.Raw, inter.Raw}, PrivateKey: ystKey, Leaf: leaf}},
		MinVersion:   tls.VersionTLS13,
		Time:         now,
	})
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return srv, roots
}

func TestYstmark(t *testing.T) {
	srv, roots := ystmark(t)
	resp, err := newClient(roots, "yst.apogee.internal", now).Get(srv.URL + "/status")
	if err != nil {
		t.Fatalf("GET /status = %v, want 200 and %q", err, "ok YST\n")
	}
	defer resp.Body.Close()
	if body, _ := io.ReadAll(resp.Body); string(body) != "ok YST\n" {
		t.Errorf("GET /status = %q, want %q", body, "ok YST\n")
	}
}

// TestImpostor plays someone who answers in Ystmark's place with a certificate for example.com, and
// counts the commands that reach them.
func TestImpostor(t *testing.T) {
	_, roots := ystmark(t)
	var received atomic.Int32
	impostor := serve(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusAccepted)
	}), nil)
	resp, err := newClient(roots, "yst.apogee.internal", now).Post(impostor.URL+"/uplink", "text/plain", strings.NewReader(sealedPing))
	if err == nil {
		resp.Body.Close()
		t.Errorf("POST /uplink to the impostor = %s, want a refused handshake", resp.Status)
	}
	if n := received.Load(); n != 0 {
		t.Errorf("commands the impostor received = %d, want 0", n)
	}
}

func TestWrongName(t *testing.T) {
	srv, roots := ystmark(t)
	resp, err := newClient(roots, "pel.apogee.internal", now).Get(srv.URL + "/status")
	if err == nil {
		resp.Body.Close()
		t.Errorf("GET /status from Ystmark, asking for pel.apogee.internal = %s, want a refused handshake", resp.Status)
	}
}

func TestTimeout(t *testing.T) {
	if got := newClient(x509.NewCertPool(), "yst.apogee.internal", now).Timeout; got != 2*time.Second {
		t.Errorf("newClient(roots, %q, now).Timeout = %v, want 2s", "yst.apogee.internal", got)
	}
}

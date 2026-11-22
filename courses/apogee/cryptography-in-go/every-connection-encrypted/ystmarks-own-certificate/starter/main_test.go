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
	"strings"
	"testing"
	"time"
)

// monday is the tests' clock: Monday 25 September 2028, 09:00 UTC.
var monday = time.Date(2028, 9, 25, 9, 0, 0, 0, time.UTC)

// now returns the tests' clock.
func now() time.Time { return monday }

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// station holds what Apogee's CA made for Ystmark: the pool holding the root alone, the
// intermediate, Ystmark's certificate 1001 and its key.
type station struct {
	roots *x509.CertPool
	inter *x509.Certificate
	leaf  *x509.Certificate
	key   ed25519.PrivateKey
}

// ystmark builds Apogee's test CA and the certificate it issued Ystmark on 22 September.
func ystmark(t *testing.T) station {
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
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return station{roots: roots, inter: inter, leaf: leaf, key: ystKey}
}

// serve starts an HTTPS server answering "ok YST" on 127.0.0.1 with cfg. Refused handshakes print nothing.
func serve(t *testing.T, cfg *tls.Config) *httptest.Server {
	t.Helper()
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok YST\n")
	}))
	srv.TLS = cfg
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv
}

// client returns a client that trusts roots alone, asks for yst.apogee.internal and reads clock.
func client(roots *x509.CertPool, clock func() time.Time) *http.Client {
	return &http.Client{Timeout: time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{
		RootCAs:    roots,
		ServerName: "yst.apogee.internal",
		MinVersion: tls.VersionTLS13,
		Time:       clock,
	}}}
}

// describe lists a chain as its serial numbers and names.
func describe(chain []*x509.Certificate) string {
	var parts []string
	for _, c := range chain {
		parts = append(parts, fmt.Sprintf("%v %s", c.SerialNumber, c.Subject.CommonName))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func TestServes(t *testing.T) {
	yst := ystmark(t)
	srv := serve(t, serverConfig(yst.leaf, yst.inter, yst.key, now))
	resp, err := client(yst.roots, now).Get(srv.URL + "/status")
	if err != nil {
		t.Fatalf("GET /status = %v, want 200 and %q", err, "ok YST\n")
	}
	defer resp.Body.Close()
	if body, _ := io.ReadAll(resp.Body); string(body) != "ok YST\n" {
		t.Errorf("GET /status = %q, want %q", body, "ok YST\n")
	}
	chain := resp.TLS.PeerCertificates
	if len(chain) != 2 || chain[0].SerialNumber.Int64() != 1001 || chain[1].Subject.CommonName != "Apogee Stations CA 2028" {
		t.Errorf("chain presented = %s, want [1001 yst.apogee.internal, 2 Apogee Stations CA 2028]", describe(chain))
	}
}

func TestExpired(t *testing.T) {
	yst := ystmark(t)
	srv := serve(t, serverConfig(yst.leaf, yst.inter, yst.key, now))
	late := yst.leaf.NotAfter.Add(time.Second)
	resp, err := client(yst.roots, func() time.Time { return late }).Get(srv.URL + "/status")
	if err == nil {
		resp.Body.Close()
		t.Errorf("GET /status from a client whose clock reads %s = %s, want a refused handshake", late.Format(time.RFC3339), resp.Status)
	}
}

func TestClock(t *testing.T) {
	yst := ystmark(t)
	cfg := serverConfig(yst.leaf, yst.inter, yst.key, now)
	if cfg.Time == nil {
		t.Fatalf("serverConfig(...).Time = nil, want now")
	}
	if got := cfg.Time(); !got.Equal(monday) {
		t.Errorf("serverConfig(...).Time() = %s, want %s", got.Format(time.RFC3339), monday.Format(time.RFC3339))
	}
}

package main

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// at returns 17 October 2028, Rotation Day, at the UTC time hh:mm.
func at(hh, mm int) time.Time { return time.Date(2028, 10, 17, hh, mm, 0, 0, time.UTC) }

// The days the stations CA issued Cape Wren's certificates: 1005 on 22 September, 1105 for Rotation Day.
var (
	sept22 = time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC)
	oct9   = time.Date(2028, 10, 9, 0, 0, 0, 0, time.UTC)
)

// clock is the test's clock, which the test moves on as Rotation Day's morning goes by.
type clock struct {
	mu sync.Mutex
	t  time.Time
}

// Now returns the time the clock was last set to.
func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

// Set moves the clock to t.
func (c *clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = t
}

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// pki is Apogee's test CA and the impostor's root of the same name.
type pki struct {
	roots    *x509.CertPool
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
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return pki{roots: roots, inter: Issuer{Cert: inter, Key: interKey}, impostor: Issuer{Cert: imp, Key: impKey}}
}

// capeWren starts Cape Wren's uplinkd presenting certificate serial, issued on day, and returns its address.
func capeWren(t *testing.T, p pki, serial int64, day time.Time, clk *clock) string {
	t.Helper()
	key := testKey(fmt.Sprintf("cert cwr.apogee.internal %d", serial))
	leaf, err := issueServer(p.inter, "cwr.apogee.internal", serial, key.Public().(ed25519.PublicKey), day)
	if err != nil {
		t.Fatalf("issueServer: %v", err)
	}
	cfg := serverConfig(leaf, p.inter.Cert, key, clk.Now)
	harden(cfg)
	requireClients(cfg, p.roots)
	srv := httptest.NewUnstartedServer(http.NotFoundHandler())
	srv.TLS = cfg
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv.Listener.Addr().String()
}

// uplinkd returns the one config uplinkd dials Cape Wren with all morning, its client certificate 2101,
// checking revocations with current.
func uplinkd(t *testing.T, p pki, current func() []byte, clk *clock) *tls.Config {
	t.Helper()
	key := testKey("cert uplinkd 2101")
	leaf, err := issueClient(p.inter, "uplinkd", 2101, key.Public().(ed25519.PublicKey), oct9)
	if err != nil {
		t.Fatalf("issueClient: %v", err)
	}
	cert := tls.Certificate{Certificate: [][]byte{leaf.Raw, p.inter.Cert.Raw}, PrivateKey: key, Leaf: leaf}
	cfg := clientConfig(p.roots, "cwr.apogee.internal", cert, clk.Now)
	withRevocations(cfg, current, p.inter.Cert, clk.Now)
	return cfg
}

// connect makes one handshake with the station at addr using cfg, and returns its error.
func connect(cfg *tls.Config, addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := (&tls.Dialer{Config: cfg}).DialContext(ctx, "tcp", addr)
	if err == nil {
		conn.Close()
	}
	return err
}

// request asks the station at addr for its status over a new connection made with cfg, and returns the
// error that stopped it, unwrapped from the request's.
func request(cfg *tls.Config, addr string) error {
	c := &http.Client{Timeout: time.Second, Transport: &http.Transport{TLSClientConfig: cfg, DisableKeepAlives: true}}
	resp, err := c.Get("https://" + addr + "/status")
	if err != nil {
		return errors.Unwrap(err)
	}
	resp.Body.Close()
	return nil
}

// list1 is revocation list 1 as the stations CA, or the impostor's root, signs it at 10:12: certificate 1005.
func list1(t *testing.T, ca Issuer) []byte {
	t.Helper()
	der, err := revoke(ca, 1, []x509.RevocationListEntry{{SerialNumber: big.NewInt(1005), RevocationTime: at(10, 12)}}, at(10, 12))
	if err != nil {
		t.Fatalf("revoke: %v", err)
	}
	return der
}

func TestNewCertificate(t *testing.T) {
	p := newPKI(t)
	clk := &clock{t: at(10, 14)}
	addr := capeWren(t, p, 1105, oct9, clk)
	list := list1(t, p.inter)
	if err := connect(uplinkd(t, p, func() []byte { return list }, clk), addr); err != nil {
		t.Errorf("uplinkd's handshake with Cape Wren presenting 1105, list 1 in force = %v, want nil", err)
	}
}

func TestRevokedCertificate(t *testing.T) {
	p := newPKI(t)
	clk := &clock{t: at(10, 14)}
	addr := capeWren(t, p, 1005, sept22, clk)
	list := list1(t, p.inter)
	if err := connect(uplinkd(t, p, func() []byte { return list }, clk), addr); err == nil || err.Error() != "certificate 1005 is revoked" {
		t.Errorf("uplinkd's handshake with Cape Wren presenting 1005, list 1 in force = %v, want %q", err, "certificate 1005 is revoked")
	}
}

// TestForgedList hands uplinkd a list the impostor signed: no certificate may be trusted on its word.
func TestForgedList(t *testing.T) {
	p := newPKI(t)
	clk := &clock{t: at(10, 14)}
	addr := capeWren(t, p, 1105, oct9, clk)
	forged := list1(t, p.impostor)
	if err := connect(uplinkd(t, p, func() []byte { return forged }, clk), addr); err == nil {
		t.Errorf("uplinkd's handshake with Cape Wren presenting 1105, a forged list in force = nil, want an error")
	}
}

// TestListChanges connects at 10:11, before any list, and again at 10:12, once list 1 is out, with the
// same config both times, as uplinkd does.
func TestListChanges(t *testing.T) {
	p := newPKI(t)
	clk := &clock{t: at(10, 11)}
	addr := capeWren(t, p, 1005, sept22, clk)
	var mu sync.Mutex
	var list []byte
	current := func() []byte {
		mu.Lock()
		defer mu.Unlock()
		return list
	}
	cfg := uplinkd(t, p, current, clk)
	if err := connect(cfg, addr); err != nil {
		t.Errorf("uplinkd's handshake at 10:11, before any list, presenting 1005 = %v, want nil", err)
	}
	signed := list1(t, p.inter)
	mu.Lock()
	list = signed
	mu.Unlock()
	clk.Set(at(10, 12))
	if err := connect(cfg, addr); err == nil || err.Error() != "certificate 1005 is revoked" {
		t.Errorf("uplinkd's handshake at 10:12, once list 1 is out, presenting 1005 = %v, want %q", err, "certificate 1005 is revoked")
	}
}

// TestResumedSession keeps the 10:11 session's ticket, so at 10:12, once list 1 is out, uplinkd resumes
// that session instead of making a new one: the check must run on the resumed handshake too.
func TestResumedSession(t *testing.T) {
	p := newPKI(t)
	clk := &clock{t: at(10, 11)}
	addr := capeWren(t, p, 1005, sept22, clk)
	var mu sync.Mutex
	var list []byte
	current := func() []byte {
		mu.Lock()
		defer mu.Unlock()
		return list
	}
	cfg := uplinkd(t, p, current, clk)
	cfg.ClientSessionCache = tls.NewLRUClientSessionCache(1)
	if err := request(cfg, addr); err != nil {
		t.Fatalf("uplinkd's request at 10:11, before any list, presenting 1005 = %v, want nil", err)
	}
	signed := list1(t, p.inter)
	mu.Lock()
	list = signed
	mu.Unlock()
	clk.Set(at(10, 12))
	if err := request(cfg, addr); err == nil || err.Error() != "certificate 1005 is revoked" {
		t.Errorf("uplinkd's request at 10:12, resuming the 10:11 session with 1005, list 1 out = %v, want %q", err, "certificate 1005 is revoked")
	}
}

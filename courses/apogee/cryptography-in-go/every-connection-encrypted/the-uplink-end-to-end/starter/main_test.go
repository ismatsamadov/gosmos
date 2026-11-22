package main

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// friday is the tests' clock: Friday 6 October 2028, 07:30 UTC.
var friday = time.Date(2028, 10, 6, 7, 30, 0, 0, time.UTC)

// now returns the tests' clock.
func now() time.Time { return friday }

// sealedPing is Cirrine-1's sealed PING from the season's test vectors, and sealKey the key that opens
// it: Cirrine-1's seal key from k-2028-07, which only the satellite holds.
const (
	sealedPing = "SEAL 2 cirrine-1 4471 k-2028-07 a1b2c3d4e5f60718293a4b5c c1f2a9798d38081d764bbcfb5e35e2598d2584a4"
	sealKey    = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"
)

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

// issue makes the certificate serial for name from ca on 22 September, a client's or a server's,
// its key made from label, and returns it with its chain.
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

// serve starts a station's uplinkd with cfg, and returns it with its queue.
func serve(t *testing.T, cfg *tls.Config) (*httptest.Server, *Queue) {
	t.Helper()
	q := &Queue{}
	srv := httptest.NewUnstartedServer(stationHandler("YST", q))
	srv.TLS = cfg
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv, q
}

// station starts a station's uplinkd with stationServer, presenting cert from ca, and returns it with
// its queue.
func station(t *testing.T, p pki, cert tls.Certificate, ca Issuer) (*httptest.Server, *Queue) {
	t.Helper()
	return serve(t, stationServer(cert.Leaf, ca.Cert, cert.PrivateKey.(ed25519.PrivateKey), p.roots, now))
}

// oldStation starts Ystmark's uplinkd as an old build of it would serve: its own certificate, clients
// proven, but nothing newer than TLS 1.2.
func oldStation(t *testing.T, p pki) (*httptest.Server, *Queue) {
	t.Helper()
	return serve(t, &tls.Config{
		Certificates: []tls.Certificate{issue(t, p.inter, "yst.apogee.internal", 1001, "cert yst.apogee.internal 1001", false)},
		MaxVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    p.roots,
		Time:         now,
	})
}

// ystmark starts Ystmark's uplinkd with its own certificate, 1001.
func ystmark(t *testing.T, p pki) (*httptest.Server, *Queue) {
	t.Helper()
	return station(t, p, issue(t, p.inter, "yst.apogee.internal", 1001, "cert yst.apogee.internal 1001", false), p.inter)
}

// brief shortens every long field of each line to its first 8 characters and an ellipsis.
func brief(lines []string) []string {
	var out []string
	for _, line := range lines {
		f := strings.Fields(line)
		for i := range f {
			if len(f[i]) > 20 {
				f[i] = f[i][:8] + "…"
			}
		}
		out = append(out, strings.Join(f, " "))
	}
	return out
}

func TestSealedLineArrives(t *testing.T) {
	p := newPKI(t)
	srv, q := ystmark(t, p)
	c := uplinkClient(p.roots, "yst.apogee.internal", issue(t, p.inter, "uplinkd", 2001, "cert uplinkd 2001", true), now)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := send(ctx, c, srv.URL, sealedPing); err != nil {
		t.Fatalf("send(ctx, uplinkClient, Ystmark, sealed PING) = %v, want nil", err)
	}
	lines := q.Lines()
	if len(lines) != 1 || lines[0] != sealedPing {
		t.Fatalf("Ystmark's queue = %q, want [%q]", brief(lines), brief([]string{sealedPing})[0])
	}
	key, _ := hex.DecodeString(sealKey)
	cirrine := &Satellite{ID: "cirrine-1", Keys: map[string]SatKeys{"k-2028-07": {Seal: key}}, Last: 4470, Now: now}
	if cmd, err := cirrine.Accept(lines[0]); err != nil || cmd.Seq != 4471 || cmd.Op != "PING" {
		t.Errorf("Cirrine-1's stand-in: Accept(the queued line) = %d %s, %v, want 4471 PING, nil", cmd.Seq, cmd.Op, err)
	}
}

func TestVersion(t *testing.T) {
	p := newPKI(t)
	srv, _ := ystmark(t, p)
	c := uplinkClient(p.roots, "yst.apogee.internal", issue(t, p.inter, "uplinkd", 2001, "cert uplinkd 2001", true), now)
	resp, err := c.Get(srv.URL + "/status")
	if err != nil {
		t.Fatalf("GET /status = %v, want 200", err)
	}
	resp.Body.Close()
	if v := resp.TLS.Version; v != tls.VersionTLS13 {
		t.Errorf("GET /status: version = %s, want TLS 1.3", tls.VersionName(v))
	}
}

func TestTimeout(t *testing.T) {
	p := newPKI(t)
	c := uplinkClient(p.roots, "yst.apogee.internal", issue(t, p.inter, "uplinkd", 2001, "cert uplinkd 2001", true), now)
	if c.Timeout != 2*time.Second {
		t.Errorf("uplinkClient(...).Timeout = %v, want 2s", c.Timeout)
	}
}

// TestRefused plays five callers or stations the uplink must turn away, each within a second.
func TestRefused(t *testing.T) {
	p := newPKI(t)
	uplinkd := issue(t, p.inter, "uplinkd", 2001, "cert uplinkd 2001", true)
	tests := []struct {
		name   string
		start  func() (*httptest.Server, *Queue)
		client func() *http.Client
	}{
		{"an impostor station", func() (*httptest.Server, *Queue) {
			return station(t, p, issue(t, p.impostor, "yst.apogee.internal", 1001, "cert impostor", false), p.impostor)
		}, func() *http.Client { return uplinkClient(p.roots, "yst.apogee.internal", uplinkd, now) }},
		{"an impostor uplinkd", func() (*httptest.Server, *Queue) { return ystmark(t, p) }, func() *http.Client {
			return uplinkClient(p.roots, "yst.apogee.internal", issue(t, p.impostor, "uplinkd", 2001, "cert impostor", true), now)
		}},
		{"uplinkd limited to TLS 1.2", func() (*httptest.Server, *Queue) { return ystmark(t, p) }, func() *http.Client {
			return &http.Client{Timeout: time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{
				RootCAs: p.roots, ServerName: "yst.apogee.internal", Certificates: []tls.Certificate{uplinkd},
				MaxVersion: tls.VersionTLS12, Time: now,
			}}}
		}},
		{"a station still on TLS 1.2", func() (*httptest.Server, *Queue) { return oldStation(t, p) }, func() *http.Client {
			return uplinkClient(p.roots, "yst.apogee.internal", uplinkd, now)
		}},
		{"passplan", func() (*httptest.Server, *Queue) { return ystmark(t, p) }, func() *http.Client {
			return uplinkClient(p.roots, "yst.apogee.internal", issue(t, p.inter, "passplan", 2002, "cert passplan 2002", true), now)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, q := tt.start()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			start := time.Now()
			err := send(ctx, tt.client(), srv.URL, sealedPing)
			if took := time.Since(start); took > time.Second {
				t.Errorf("send with %s took %v, want under a second", tt.name, took.Round(time.Millisecond))
			}
			if err == nil {
				t.Errorf("send(ctx, %s, sealed PING) = nil, want a refusal", tt.name)
			}
			if lines := q.Lines(); len(lines) != 0 {
				t.Errorf("the station's queue after %s = %q, want it empty", tt.name, brief(lines))
			}
		})
	}
}

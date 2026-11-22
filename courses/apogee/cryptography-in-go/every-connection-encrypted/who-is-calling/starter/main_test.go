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

// thursday is the tests' clock: Thursday 28 September 2028, 09:00 UTC.
var thursday = time.Date(2028, 9, 28, 9, 0, 0, 0, time.UTC)

// now returns the tests' clock.
func now() time.Time { return thursday }

// Two lines uplinkd sends Ystmark: Cirrine-1's sealed PING from the season's test vectors, and
// Pipit-2's signed command from Season 1's lab bench.
const (
	sealedPing = "SEAL 2 cirrine-1 4471 k-2028-07 a1b2c3d4e5f60718293a4b5c c1f2a9798d38081d764bbcfb5e35e2598d2584a4"
	signedCmd  = "CMD pipit-2 118 SET_MODE SAFE sig=512545bcb43a5bb6ccb9f9c8511f4cc8c5fedeb38c226376af5b5aaae60e3b4d"
)

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// apogee is Apogee's test CA and what it issued on 22 September.
type apogee struct {
	roots    *x509.CertPool
	root     *x509.Certificate
	inter    Issuer
	yst      tls.Certificate // yst.apogee.internal, 1001
	uplinkd  tls.Certificate // 2001
	passplan tls.Certificate // 2002
}

// newApogee builds the CA and its certificates from the season's labels.
func newApogee(t *testing.T) apogee {
	t.Helper()
	caDay := time.Date(2028, 9, 18, 0, 0, 0, 0, time.UTC)
	rootKey, interKey := testKey("ca root 2028"), testKey("ca stations 2028")
	root, err := newRoot(rootKey, caDay)
	if err != nil {
		t.Fatalf("newRoot: %v", err)
	}
	inter, err := newIntermediate(Issuer{Cert: root, Key: rootKey}, interKey, caDay)
	if err != nil {
		t.Fatalf("newIntermediate: %v", err)
	}
	a := apogee{roots: x509.NewCertPool(), root: root, inter: Issuer{Cert: inter, Key: interKey}}
	a.roots.AddCert(root)
	a.yst = a.issue(t, issueServer, "yst.apogee.internal", 1001)
	a.uplinkd = a.issue(t, issueClient, "uplinkd", 2001)
	a.passplan = a.issue(t, issueClient, "passplan", 2002)
	return a
}

// issue makes the certificate serial for name with mint, its key made from "cert <name> <serial>",
// and returns it with its chain.
func (a apogee) issue(t *testing.T, mint func(Issuer, string, int64, ed25519.PublicKey, time.Time) (*x509.Certificate, error), name string, serial int64) tls.Certificate {
	t.Helper()
	key := testKey(fmt.Sprintf("cert %s %d", name, serial))
	leaf, err := mint(a.inter, name, serial, key.Public().(ed25519.PublicKey), time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("issuing %s %d: %v", name, serial, err)
	}
	return tls.Certificate{Certificate: [][]byte{leaf.Raw, a.inter.Cert.Raw}, PrivateKey: key, Leaf: leaf}
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

// ystmark starts Ystmark's uplinkd: both ends proven, its handler stationHandler, its queue q.
func ystmark(t *testing.T, a apogee, q *Queue) *httptest.Server {
	t.Helper()
	cfg := serverConfig(a.yst.Leaf, a.inter.Cert, a.yst.PrivateKey.(ed25519.PrivateKey), now)
	harden(cfg)
	requireClients(cfg, a.roots)
	srv := httptest.NewUnstartedServer(stationHandler("YST", q))
	srv.TLS = cfg
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv
}

// do sends method and path to the server at url as the holder of cert, and returns the answer's
// status and body.
func do(t *testing.T, a apogee, cert tls.Certificate, url, method, path, body string) (int, string) {
	t.Helper()
	c := &http.Client{Timeout: time.Second, Transport: &http.Transport{
		TLSClientConfig: clientConfig(a.roots, "yst.apogee.internal", cert, now),
	}}
	req, err := http.NewRequest(method, url+path, strings.NewReader(body))
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("%s %s = %v, want an answer", method, path, err)
	}
	defer resp.Body.Close()
	answer, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(answer)
}

func TestUplinkdQueues(t *testing.T) {
	a := newApogee(t)
	q := &Queue{}
	srv := ystmark(t, a, q)
	for _, tt := range []struct{ line, want string }{
		{sealedPing, "queued cirrine-1 4471\n"},
		{signedCmd, "queued pipit-2 118\n"},
	} {
		if code, body := do(t, a, a.uplinkd, srv.URL, "POST", "/uplink", tt.line); code != http.StatusAccepted || body != tt.want {
			t.Errorf("POST /uplink as uplinkd = %d %q, want 202 %q", code, body, tt.want)
		}
	}
	if got := q.Lines(); len(got) != 2 || got[0] != sealedPing || got[1] != signedCmd {
		t.Errorf("Ystmark's queue = %q, want both lines, in order", brief(got))
	}
}

func TestPassplan(t *testing.T) {
	a := newApogee(t)
	q := &Queue{}
	srv := ystmark(t, a, q)
	if code, body := do(t, a, a.passplan, srv.URL, "POST", "/uplink", sealedPing); code != http.StatusForbidden || body != "only uplinkd may queue commands\n" {
		t.Errorf("POST /uplink as passplan = %d %q, want 403 %q", code, body, "only uplinkd may queue commands\n")
	}
	if got := q.Lines(); len(got) != 0 {
		t.Errorf("Ystmark's queue after passplan's POST = %q, want it empty", brief(got))
	}
	if code, body := do(t, a, a.passplan, srv.URL, "GET", "/status", ""); code != http.StatusOK || body != "ok YST\n" {
		t.Errorf("GET /status as passplan = %d %q, want 200 %q", code, body, "ok YST\n")
	}
}

// callerOf runs caller, turning a panic into a failure.
func callerOf(t *testing.T, r *http.Request) (name string, ok bool) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Errorf("caller(request) panicked: %v, want an answer", p)
		}
	}()
	return caller(r)
}

func TestCaller(t *testing.T) {
	a := newApogee(t)
	tests := []struct {
		name   string
		state  *tls.ConnectionState
		want   string
		wantOK bool
	}{
		{"a request without TLS", nil, "", false},
		{"TLS without a verified chain", &tls.ConnectionState{}, "", false},
		{"uplinkd's certificate, sent but never verified", &tls.ConnectionState{PeerCertificates: []*x509.Certificate{a.uplinkd.Leaf}}, "", false},
		{"uplinkd's verified chain", &tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{a.uplinkd.Leaf},
			VerifiedChains:   [][]*x509.Certificate{{a.uplinkd.Leaf, a.inter.Cert, a.root}},
		}, "uplinkd", true},
		{"passplan's verified chain", &tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{a.passplan.Leaf},
			VerifiedChains:   [][]*x509.Certificate{{a.passplan.Leaf, a.inter.Cert, a.root}},
		}, "passplan", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/uplink", nil)
			r.TLS = tt.state
			if got, ok := callerOf(t, r); got != tt.want || ok != tt.wantOK {
				t.Errorf("caller(%s) = %q, %v, want %q, %v", tt.name, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

// TestWithoutTLS hands the handler a request that never went through a handshake.
func TestWithoutTLS(t *testing.T) {
	q := &Queue{}
	rec := httptest.NewRecorder()
	func() {
		defer func() {
			if p := recover(); p != nil {
				t.Errorf("POST /uplink without TLS panicked: %v, want 403", p)
			}
		}()
		stationHandler("YST", q).ServeHTTP(rec, httptest.NewRequest("POST", "/uplink", strings.NewReader(sealedPing)))
	}()
	if rec.Code != http.StatusForbidden {
		t.Errorf("POST /uplink without TLS = %d %q, want 403 %q", rec.Code, rec.Body.String(), "only uplinkd may queue commands\n")
	}
	if got := q.Lines(); len(got) != 0 {
		t.Errorf("Ystmark's queue after a POST without TLS = %q, want it empty", brief(got))
	}
}

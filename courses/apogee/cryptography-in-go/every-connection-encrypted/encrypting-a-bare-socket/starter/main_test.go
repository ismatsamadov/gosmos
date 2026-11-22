package main

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

// tuesday is the tests' clock: Tuesday 3 October 2028, 09:00 UTC.
var tuesday = time.Date(2028, 10, 3, 9, 0, 0, 0, time.UTC)

// now returns the tests' clock.
func now() time.Time { return tuesday }

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

// controller starts a stand-in dish controller on 127.0.0.1 behind TLS, presenting cert from ca and
// asking for a client certificate: it echoes each line it reads, PING yst in, PING yst out. It
// returns the controller's address; the test's cleanup stops it.
func controller(t *testing.T, p pki, cert tls.Certificate, ca Issuer) string {
	t.Helper()
	cfg := serverConfig(cert.Leaf, ca.Cert, cert.PrivateKey.(ed25519.PrivateKey), now)
	harden(cfg)
	requireClients(cfg, p.roots)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	var wg sync.WaitGroup
	t.Cleanup(func() {
		ln.Close()
		wg.Wait()
	})
	tl := tls.NewListener(ln, cfg)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			c, err := tl.Accept()
			if err != nil {
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer c.Close()
				c.SetDeadline(time.Now().Add(time.Second))
				lines := bufio.NewScanner(c)
				for lines.Scan() {
					fmt.Fprintf(c, "%s\n", lines.Text())
				}
			}()
		}
	}()
	return ln.Addr().String()
}

// dishd returns the config dishd dials the controllers with: Apogee's root, Ystmark's name and
// dishd's client certificate, 2003.
func dishd(t *testing.T, p pki) *tls.Config {
	t.Helper()
	return clientConfig(p.roots, "yst.apogee.internal", issue(t, p.inter, "dishd", 2003, "cert dishd 2003", true), now)
}

func TestEcho(t *testing.T) {
	p := newPKI(t)
	addr := controller(t, p, issue(t, p.inter, "yst.apogee.internal", 1001, "cert yst.apogee.internal 1001", false), p.inter)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := dialController(ctx, addr, dishd(t, p))
	if err != nil {
		t.Fatalf("dialController(ctx, controller, cfg) = %v, want a connection", err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := fmt.Fprint(conn, "PING yst\n"); err != nil {
		t.Fatalf("writing PING yst = %v, want nil", err)
	}
	if line, err := bufio.NewReader(conn).ReadString('\n'); line != "PING yst\n" {
		t.Errorf("the controller's answer to %q = %q, %v, want %q", "PING yst\n", line, err, "PING yst\n")
	}
}

// TestSilentController dials a listener that takes connections and never says a word. A dial that
// ignores ctx would wait forever, so the test stops waiting after 2 seconds.
func TestSilentController(t *testing.T) {
	p := newPKI(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	defer ln.Close()
	cfg := dishd(t, p)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	type dialed struct {
		conn net.Conn
		err  error
	}
	done := make(chan dialed, 1)
	start := time.Now()
	go func() {
		conn, err := dialController(ctx, ln.Addr().String(), cfg)
		done <- dialed{conn, err}
	}()
	select {
	case d := <-done:
		took := time.Since(start)
		if d.err == nil {
			d.conn.Close()
			t.Errorf("dialController to a listener that never answers = a connection after %v, want an error within 200ms", took.Round(time.Millisecond))
		} else if took > time.Second {
			t.Errorf("dialController to a listener that never answers took %v, want an error within 200ms", took.Round(time.Millisecond))
		}
	case <-time.After(2 * time.Second):
		t.Errorf("dialController to a listener that never answers was still waiting after 2s, want an error within 200ms: the handshake ignores ctx")
	}
}

func TestImpostorController(t *testing.T) {
	p := newPKI(t)
	addr := controller(t, p, issue(t, p.impostor, "yst.apogee.internal", 1001, "cert impostor", false), p.impostor)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := dialController(ctx, addr, dishd(t, p))
	if err == nil {
		conn.Close()
		t.Errorf("dialController to a controller with the impostor's certificate = a connection, want an error")
	}
}

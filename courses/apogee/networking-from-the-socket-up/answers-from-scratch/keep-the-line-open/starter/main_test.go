package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestServeConn(t *testing.T) {
	t.Run("the collector's three requests", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536}}
		s.Store.Put("furrow-a/YST/18856", testRecording(4096))
		ln := startServer(t, s)
		client := newClient(t, ln.Addr().String())
		for _, r := range []struct{ method, path string }{
			{"HEAD", "/recordings/furrow-a/YST/18856"},
			{"GET", "/recordings/furrow-a/YST/18856"},
			{"GET", "/recordings"},
		} {
			req, err := http.NewRequest(r.method, "http://127.0.0.1"+r.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("%s %s failed: %v", r.method, r.path, err)
			}
			io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode != 200 {
				t.Errorf("%s %s = %d, want 200", r.method, r.path, res.StatusCode)
			}
		}
		if n := ln.accepted.Load(); n != 1 {
			t.Errorf("HEAD, GET and GET /recordings took %d connections, want 1", n)
		}
	})
	t.Run("a request with Connection: close", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536}}
		res, err := exchange(t, startServer(t, s).Addr().String(), "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\nConnection: close\r\n\r\n")
		if res.StatusCode != 200 || !res.Close || err != nil {
			t.Errorf("GET /recordings with Connection: close = %d (Connection: close %v), then %v; want 200 (Connection: close true), then the end of the stream", res.StatusCode, res.Close, err)
		}
	})
	t.Run("a 404", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536}}
		res, err := exchange(t, startServer(t, s).Addr().String(), "GET /recordings/furrow-a/YST/18857 HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
		if res.StatusCode != 404 || !res.Close || err != nil {
			t.Errorf("GET of a pass not kept = %d (Connection: close %v), then %v; want 404 (Connection: close true), then the end of the stream", res.StatusCode, res.Close, err)
		}
	})
	t.Run("a connection quiet after its answer, Idle 50ms", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536, Idle: 50 * time.Millisecond}}
		res, err := exchange(t, startServer(t, s).Addr().String(), "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
		if res.StatusCode != 200 || res.Close || err != nil {
			t.Errorf("GET /recordings = %d (Connection: close %v), then %v; want 200 (Connection: close false), then the end of the stream once Idle has passed", res.StatusCode, res.Close, err)
		}
	})
}

// exchange sends raw to the server at addr on a connection of its own, reads the answer,
// then waits up to a second for the server to hang up. It returns the answer and what that
// wait ended with: nil when the stream ended.
func exchange(t *testing.T, addr, raw string) (*http.Response, error) {
	t.Helper()
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.WriteString(conn, raw); err != nil {
		t.Fatalf("sending the request: %v", err)
	}
	br := bufio.NewReader(conn)
	res, err := http.ReadResponse(br, nil)
	if err != nil {
		t.Fatalf("reading the answer: %v", err)
	}
	io.ReadAll(res.Body)
	res.Body.Close()
	if _, err := br.ReadByte(); !errors.Is(err, io.EOF) {
		if err == nil {
			err = errors.New("more bytes after the answer")
		}
		return res, err
	}
	return res, nil
}

// countingListener counts the connections it accepts.
type countingListener struct {
	net.Listener
	accepted atomic.Int32
}

// Accept accepts a connection and counts it.
func (l *countingListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err == nil {
		l.accepted.Add(1)
	}
	return c, err
}

// testRecording is n bytes of a test recording: byte i is i % 256, so it holds zeros, CRs
// and LFs, a newline every 256 bytes.
func testRecording(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 256)
	}
	return b
}

// june19 is the station server's clock in these tests.
func june19() time.Time { return time.Date(2028, 6, 19, 15, 30, 0, 0, time.UTC) }

// startServer serves s on a loopback port the system picks and returns its listener. When
// the test ends, it closes the listener and waits until every connection has been served.
func startServer(t *testing.T, s *Server) *countingListener {
	t.Helper()
	inner, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ln := &countingListener{Listener: inner}
	served := make(chan error, 1)
	go func() { served <- s.Serve(ln) }()
	t.Cleanup(func() {
		ln.Close()
		select {
		case <-served:
		case <-time.After(time.Second):
			t.Error("the station server was still serving a second after its listener closed")
		}
	})
	return ln
}

// newClient returns an HTTP client, like the collector's, whose every connection goes to
// addr, so its requests can name http://127.0.0.1/.
func newClient(t *testing.T, addr string) *http.Client {
	tr := &http.Transport{DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "tcp", addr)
	}}
	t.Cleanup(tr.CloseIdleConnections)
	return &http.Client{Timeout: time.Second, Transport: tr}
}

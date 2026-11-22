package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestReadChunkedLimit(t *testing.T) {
	recording := testRecording(65536)
	tests := []struct {
		name, stream string
		max          int64
	}{
		{"one byte past 65536, the chunk unsent", "8000\r\n" + string(recording[:32768]) + "\r\n8000\r\n" + string(recording[32768:]) + "\r\n1\r\n", 65536},
		{"one chunk of 65537 bytes, unsent", "10001\r\n", 65536},
		{"a limit of 24 bytes", "c\r\nfirst piece \r\nd\r\n", 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := readChunked(bufio.NewReaderSize(strings.NewReader(tt.stream), 1024), tt.max)
			want := fmt.Sprintf("body over %d bytes", tt.max)
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != 413 || e.Msg != want {
				t.Errorf("readChunked(%s) error = %v, want 413 %s, before reading the chunk", tt.name, err, want)
			}
		})
	}
	t.Run("exactly 65536 bytes", func(t *testing.T) {
		stream := "8000\r\n" + string(recording[:32768]) + "\r\n8000\r\n" + string(recording[32768:]) + "\r\n0\r\n\r\n"
		body, _, err := readChunked(bufio.NewReaderSize(strings.NewReader(stream), 1024), 65536)
		if !bytes.Equal(body, recording) || err != nil {
			t.Errorf("readChunked(65536 bytes in two chunks) = %s, %v, want %s, nil", describe(body), err, describe(recording))
		}
	})
	t.Run("no limit", func(t *testing.T) {
		body, _, err := readChunked(bufio.NewReaderSize(strings.NewReader("c\r\nfirst piece \r\n0\r\n\r\n"), 1024), 0)
		if string(body) != "first piece " || err != nil {
			t.Errorf("readChunked(12 bytes, max 0) = %q, %v, want %q, nil", body, err, "first piece ")
		}
	})
}

func TestTooBigThroughTheServer(t *testing.T) {
	s := &Server{Store: NewStore(), Now: june8, Limits: Limits{MaxBody: 65536, Read: 50 * time.Millisecond}}
	addr := startServer(t, s)
	recording := testRecording(65536)

	t.Run("65536 bytes, chunked by Go's client", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/furrow-b/YST/18672", io.MultiReader(bytes.NewReader(recording)))
		if err != nil {
			t.Fatal(err)
		}
		res, err := newClient(t, addr).Do(req)
		if err != nil {
			t.Fatalf("PUT of 65536 bytes failed: %v", err)
		}
		res.Body.Close()
		got, _ := s.Store.Get("furrow-b/YST/18672")
		if res.StatusCode != 201 || !bytes.Equal(got, recording) {
			t.Errorf("PUT of 65536 bytes = %d, keeping %s, want 201, keeping %s", res.StatusCode, describe(got), describe(recording))
		}
	})
	t.Run("one byte more, its chunk unsent", func(t *testing.T) {
		raw := "PUT /recordings/furrow-a/YST/18674 HTTP/1.1\r\nHost: 127.0.0.1\r\nTransfer-Encoding: chunked\r\n\r\n" +
			"10000\r\n" + string(recording) + "\r\n1\r\n"
		res, rest := roundTrip(t, addr, raw)
		wantAnswer(t, "PUT of 65537 bytes, chunked", res, rest, 413, "body over 65536 bytes\n")
	})
	t.Run("a Content-Length of 65537, no body sent", func(t *testing.T) {
		raw := "PUT /recordings/furrow-a/YST/18674 HTTP/1.1\r\nHost: 127.0.0.1\r\nContent-Length: 65537\r\n\r\n"
		res, rest := roundTrip(t, addr, raw)
		wantAnswer(t, "PUT with Content-Length 65537", res, rest, 413, "body over 65536 bytes\n")
	})
	if _, ok := s.Store.Get("furrow-a/YST/18674"); ok {
		t.Error("the store keeps furrow-a/YST/18674, want nothing")
	}
}

// wantAnswer fails the test unless res has status and body, said Connection: close, and
// nothing came after it.
func wantAnswer(t *testing.T, call string, res *http.Response, rest string, status int, body string) {
	t.Helper()
	got, _ := io.ReadAll(res.Body)
	if res.StatusCode != status || string(got) != body || !res.Close {
		t.Errorf("%s = %d %q (Connection: close %v), want %d %q (Connection: close true)", call, res.StatusCode, got, res.Close, status, body)
	}
	if rest != "" {
		t.Errorf("%s: after the answer the server sent %q, want nothing more", call, rest)
	}
}

// roundTrip sends raw to the server at addr on a connection of its own and reads its answer,
// then everything the server sends after it until it hangs up.
func roundTrip(t *testing.T, addr, raw string) (*http.Response, string) {
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
	body, _ := io.ReadAll(res.Body)
	res.Body = io.NopCloser(bytes.NewReader(body))
	rest, err := io.ReadAll(br)
	if err != nil {
		t.Errorf("after the answer, the connection gave %v, want the server to hang up", err)
	}
	return res, string(rest)
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

// describe sums up a body for a failure message: its size and its CRC-32.
func describe(b []byte) string {
	return fmt.Sprintf("%d bytes (crc %08x)", len(b), crc32.ChecksumIEEE(b))
}

// june8 is the station server's clock in these tests.
func june8() time.Time { return time.Date(2028, 6, 8, 8, 30, 0, 0, time.UTC) }

// startServer serves s on a loopback port the system picks and returns its address. When
// the test ends, it closes the listener and waits until every connection has been served.
func startServer(t *testing.T, s *Server) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
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
	return ln.Addr().String()
}

// newClient returns an HTTP client, like the recorder's, whose every connection goes to
// addr, so its requests can name http://127.0.0.1/.
func newClient(t *testing.T, addr string) *http.Client {
	tr := &http.Transport{DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "tcp", addr)
	}}
	t.Cleanup(tr.CloseIdleConnections)
	return &http.Client{Timeout: time.Second, Transport: tr}
}

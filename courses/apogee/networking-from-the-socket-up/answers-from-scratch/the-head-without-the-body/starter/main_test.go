package main

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRecording(t *testing.T) {
	s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536}}
	recording := testRecording(4096)
	s.Store.Put("furrow-a/YST/18856", recording)
	addr := startServer(t, s)
	client := newClient(t, addr)
	const url = "http://127.0.0.1/recordings/furrow-a/YST/18856"

	t.Run("GET", func(t *testing.T) {
		res, err := client.Get(url)
		if err != nil {
			t.Fatalf("GET /recordings/furrow-a/YST/18856 failed: %v", err)
		}
		got, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode != 200 || !bytes.Equal(got, recording) {
			t.Errorf("GET /recordings/furrow-a/YST/18856 = %d, %s, want 200, %s", res.StatusCode, describe(got), describe(recording))
		}
		wantFields(t, "GET /recordings/furrow-a/YST/18856", res.Header, "application/octet-stream", "a2912082")
	})
	t.Run("HEAD through Go's client", func(t *testing.T) {
		res, err := client.Head(url)
		if err != nil {
			t.Fatalf("HEAD /recordings/furrow-a/YST/18856 failed: %v", err)
		}
		res.Body.Close()
		if res.StatusCode != 200 || res.ContentLength != 4096 {
			t.Errorf("HEAD /recordings/furrow-a/YST/18856 = %d with ContentLength %d, want 200 with 4096, the GET's length", res.StatusCode, res.ContentLength)
		}
		wantFields(t, "HEAD /recordings/furrow-a/YST/18856", res.Header, "application/octet-stream", "a2912082")
	})
	t.Run("HEAD on a raw connection", func(t *testing.T) {
		head, body := exchange(t, addr, "HEAD /recordings/furrow-a/YST/18856 HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
		if !strings.Contains(head, "\r\nContent-Length: 4096\r\n") || body != "" {
			t.Errorf("HEAD /recordings/furrow-a/YST/18856 sent the head\n%q\nthen %d bytes after the blank line; want Content-Length: 4096, then nothing", head, len(body))
		}
	})
	for _, method := range []string{"GET", "HEAD"} {
		t.Run(method+" of a pass not kept", func(t *testing.T) {
			head, body := exchange(t, addr, method+" /recordings/furrow-a/YST/18857 HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
			want := "no recording furrow-a/YST/18857\n"
			if method == "HEAD" {
				want = ""
			}
			if !strings.HasPrefix(head, "HTTP/1.1 404 Not Found\r\n") || body != want {
				t.Errorf("%s /recordings/furrow-a/YST/18857 = %q, then %q; want a 404, then %q", method, firstLine(head), body, want)
			}
		})
	}
}

// wantFields fails the test unless h carries the recording's type and CRC.
func wantFields(t *testing.T, call string, h http.Header, contentType, crc string) {
	t.Helper()
	if got := h.Get("Content-Type"); got != contentType {
		t.Errorf("%s answered Content-Type %q, want %q", call, got, contentType)
	}
	if got := h.Get("Apogee-Crc32"); got != crc {
		t.Errorf("%s answered Apogee-Crc32 %q, want %q", call, got, crc)
	}
}

// exchange sends raw to the server at addr on a connection of its own and returns the head
// of the answer, up to its blank line, and everything after it until the server hangs up.
func exchange(t *testing.T, addr, raw string) (head, rest string) {
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
	all, err := io.ReadAll(conn)
	if err != nil {
		t.Errorf("reading the answer: %v, want the server to hang up after it", err)
	}
	head, rest, _ = strings.Cut(string(all), "\r\n\r\n")
	return head + "\r\n\r\n", rest
}

// firstLine is the status line of an answer's head.
func firstLine(head string) string {
	line, _, _ := strings.Cut(head, "\r\n")
	return line
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

// june19 is the station server's clock in these tests.
func june19() time.Time { return time.Date(2028, 6, 19, 14, 2, 0, 0, time.UTC) }

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

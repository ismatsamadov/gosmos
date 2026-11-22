package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestWriteChunk(t *testing.T) {
	tests := []struct {
		name  string
		piece []byte
		want  string
	}{
		{"4,096 bytes", testRecording(4096), "1000\r\n" + string(testRecording(4096)) + "\r\n"},
		{"one byte", []byte("Y"), "1\r\nY\r\n"},
		{"a line of the list", []byte("furrow-a/YST/18856 4096 a2912082\n"), "21\r\nfurrow-a/YST/18856 4096 a2912082\n\r\n"},
		{"nothing", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sent bytes.Buffer
			err := writeChunk(bufio.NewWriter(&sent), tt.piece)
			if got := sent.String(); got != tt.want || err != nil {
				t.Errorf("writeChunk(%s) sent %q, %v, want %q, nil (sent without waiting for a Flush)", tt.name, clip(got), err, clip(tt.want))
			}
		})
	}
}

func TestEndChunks(t *testing.T) {
	var sent bytes.Buffer
	if err := endChunks(bufio.NewWriter(&sent)); sent.String() != "0\r\n\r\n" || err != nil {
		t.Errorf("endChunks sent %q, %v, want %q, nil (sent without waiting for a Flush)", sent.String(), err, "0\r\n\r\n")
	}
}

func TestFeed(t *testing.T) {
	t.Run("through Go's client", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536}}
		res, err := newClient(t, startServer(t, s)).Get("http://127.0.0.1/feed")
		if err != nil {
			t.Fatalf("GET /feed failed: %v", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 || len(res.TransferEncoding) != 1 || res.TransferEncoding[0] != "chunked" || res.ContentLength != -1 {
			t.Errorf("GET /feed = %d, TransferEncoding %q, ContentLength %d, want 200, [chunked], -1", res.StatusCode, res.TransferEncoding, res.ContentLength)
		}
		lines := bufio.NewReader(res.Body)
		for _, r := range []struct {
			id   string
			size int
			want string
		}{
			{"furrow-a/YST/18856", 4096, "furrow-a/YST/18856 4096 a2912082\n"},
			{"cirrine-1/YST/14302", 2048, "cirrine-1/YST/14302 2048 9f5edd58\n"},
		} {
			s.Store.Put(r.id, testRecording(r.size))
			got, err := lines.ReadString('\n')
			if got != r.want || err != nil {
				t.Fatalf("after %s was stored, the feed gave %q, %v, want %q, nil", r.id, got, err, r.want)
			}
		}
		s.Store.Close()
		if rest, err := io.ReadAll(lines); len(rest) != 0 || err != nil {
			t.Errorf("after the store closed, the feed gave %q, %v, want its end: nothing more and no error", rest, err)
		}
	})
	t.Run("a collector that leaves", func(t *testing.T) {
		gone := errors.New("broken pipe")
		lines := make(chan string, 1)
		lines <- "furrow-a/YST/18856 4096 a2912082"
		close(lines)
		s := &Server{Store: NewStore()}
		if err := s.feed(bufio.NewWriter(failingWriter{gone}), lines); !errors.Is(err, gone) {
			t.Errorf("feed to a collector that has left = %v, want the write's error, %v", err, gone)
		}
	})
}

// clip shortens a long string for a failure message, keeping its start and its end.
func clip(s string) string {
	if len(s) <= 40 {
		return s
	}
	return s[:20] + "..." + s[len(s)-12:]
}

// failingWriter is a connection whose other end has gone: every write fails with err.
type failingWriter struct{ err error }

// Write fails.
func (w failingWriter) Write(p []byte) (int, error) { return 0, w.err }

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
func june19() time.Time { return time.Date(2028, 6, 19, 17, 10, 0, 0, time.UTC) }

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
		s.Store.Close() // a feed still open ends, so its connection does
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

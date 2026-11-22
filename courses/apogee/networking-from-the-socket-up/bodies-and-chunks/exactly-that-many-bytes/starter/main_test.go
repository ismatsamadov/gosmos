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

func TestReadBody(t *testing.T) {
	t.Run("a recording sent with its head", func(t *testing.T) {
		head := "PUT /recordings/pipit-2/YST/16592 HTTP/1.1\r\nHost: 127.0.0.1\r\nContent-Length: 4096\r\n\r\n"
		br := bufio.NewReaderSize(strings.NewReader(head+string(testRecording(4096))), 1024)
		req, err := readRequest(br, 65536)
		if err != nil {
			t.Fatalf("readRequest(the upload's head) = %v; the test's head should read", err)
		}
		got, err := readBody(br, req.Length)
		if want := testRecording(4096); !bytes.Equal(got, want) || err != nil {
			t.Errorf("readBody(br, 4096) = %s, %v, want %s, nil", describe(got), err, describe(want))
		}
	})
	t.Run("an empty body", func(t *testing.T) {
		r := &countingReader{r: strings.NewReader("GET /recordings HTTP/1.1\r\n")}
		got, err := readBody(bufio.NewReader(r), 0)
		if len(got) != 0 || err != nil || r.reads != 0 {
			t.Errorf("readBody(br, 0) = %s, %v after %d reads, want 0 bytes, nil after 0 reads", describe(got), err, r.reads)
		}
	})
	t.Run("a body cut after 812 bytes", func(t *testing.T) {
		br := bufio.NewReaderSize(bytes.NewReader(testRecording(812)), 1024)
		_, err := readBody(br, 4096)
		wantRefusal(t, "readBody(812 of 4096 bytes)", err, 400, "body ended after 812 of 4096 bytes")
	})
	t.Run("an empty stream", func(t *testing.T) {
		_, err := readBody(bufio.NewReader(strings.NewReader("")), 24)
		wantRefusal(t, "readBody(0 of 24 bytes)", err, 400, "body ended after 0 of 24 bytes")
	})
	t.Run("Go's client through the station server", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june5, Limits: Limits{MaxBody: 65536, Read: 50 * time.Millisecond}}
		addr := startServer(t, s)
		req, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/pipit-2/YST/16592", bytes.NewReader(testRecording(1024)))
		if err != nil {
			t.Fatal(err)
		}
		res, err := newClient(t, addr).Do(req)
		if err != nil {
			t.Fatalf("PUT /recordings/pipit-2/YST/16592 failed: %v", err)
		}
		res.Body.Close()
		if res.StatusCode != 201 {
			t.Errorf("PUT /recordings/pipit-2/YST/16592 = %d, want 201", res.StatusCode)
		}
		got, _ := s.Store.Get("pipit-2/YST/16592")
		if want := testRecording(1024); !bytes.Equal(got, want) {
			t.Errorf("the store keeps %s, want %s", describe(got), describe(want))
		}
	})
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

// countingReader counts the reads made on it.
type countingReader struct {
	r     io.Reader
	reads int
}

// Read reads from r and counts the read.
func (c *countingReader) Read(p []byte) (int, error) {
	c.reads++
	return c.r.Read(p)
}

// wantRefusal fails the test unless err is a *RequestError with status and msg.
func wantRefusal(t *testing.T, call string, err error, status int, msg string) {
	t.Helper()
	if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != status || e.Msg != msg {
		t.Errorf("%s error = %v, want %d %s", call, err, status, msg)
	}
}

// june5 is the station server's clock in these tests.
func june5() time.Time { return time.Date(2028, 6, 5, 9, 30, 0, 0, time.UTC) }

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

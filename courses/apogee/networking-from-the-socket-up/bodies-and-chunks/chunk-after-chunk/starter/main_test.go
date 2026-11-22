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

// goUpload is what Go's client sends for a chunked PUT of 24 bytes in two pieces, with
// Apogee-Crc32 in its trailer.
const goUpload = "PUT /recordings/furrow-a/YST/18856 HTTP/1.1\r\n" +
	"Host: 127.0.0.1\r\n" +
	"User-Agent: Go-http-client/1.1\r\n" +
	"Transfer-Encoding: chunked\r\n" +
	"Trailer: Apogee-Crc32\r\n" +
	"Accept-Encoding: gzip\r\n" +
	"\r\n" +
	"c\r\nfirst piece \r\n" +
	"c\r\nsecond piece\r\n" +
	"0\r\n" +
	"Apogee-Crc32: e2fc055d\r\n" +
	"\r\n"

func TestReadChunked(t *testing.T) {
	t.Run("Go's own two-chunk upload", func(t *testing.T) {
		br := bufio.NewReaderSize(strings.NewReader(goUpload), 1024)
		req, err := readRequest(br, 65536)
		if err != nil || req.Length != -1 {
			t.Fatalf("readRequest(Go's upload) = %v; the test's head should read as chunked", err)
		}
		body, trailer, err := readChunked(br, 65536)
		if string(body) != "first piece second piece" || err != nil {
			t.Errorf("readChunked(Go's upload) = %q, %v, want %q, nil", body, err, "first piece second piece")
		}
		if got := trailer.Get("apogee-crc32"); got != "e2fc055d" {
			t.Errorf("readChunked(Go's upload) gave a trailer with apogee-crc32 %q, want %q", got, "e2fc055d")
		}
	})
	recording := testRecording(4096)
	tests := []struct {
		name, stream string
		want         []byte
	}{
		{"chunks of 1, 7 and 4088 bytes", chunked(recording, 1, 7, 4088), recording},
		{"a lone zero chunk", "0\r\n\r\n", nil},
		{"a chunk of CRs and LFs", "4\r\n\r\n\r\n\r\n0\r\n\r\n", []byte("\r\n\r\n")},
		{"an extension and a chunk of 4,096 bytes", "1000;src=furrow-b\r\n" + string(recording) + "\r\n0\r\n\r\n", recording},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, trailer, err := readChunked(bufio.NewReaderSize(strings.NewReader(tt.stream), 1024), 65536)
			if !bytes.Equal(body, tt.want) || err != nil {
				t.Errorf("readChunked(%s) = %s, %v, want %s, nil", tt.name, describe(body), err, describe(tt.want))
			}
			if len(trailer) != 0 {
				t.Errorf("readChunked(%s) gave a trailer of %d fields, want none", tt.name, len(trailer))
			}
		})
	}
	refusals := []struct {
		name, stream, msg string
	}{
		{"xx after a chunk's data", "c\r\nfirst piece xx\r\n0\r\n\r\n", "bad chunk ending"},
		{"a size line of 1,100 bytes", strings.Repeat("0", 1099) + "c\r\nfirst piece \r\n0\r\n\r\n", "chunk line over 1024 bytes"},
		{"a size that isn't hex", "zz\r\nfirst piece \r\n0\r\n\r\n", `bad chunk size "zz"`},
	}
	for _, tt := range refusals {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := readChunked(bufio.NewReaderSize(strings.NewReader(tt.stream), 1024), 65536)
			wantRefusal(t, "readChunked("+tt.name+")", err, 400, tt.msg)
		})
	}
	t.Run("Go's client through the station server", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june5, Limits: Limits{MaxBody: 65536, Read: 50 * time.Millisecond}}
		req, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/furrow-a/YST/18856", &pieces{"first piece ", "second piece"})
		if err != nil {
			t.Fatal(err)
		}
		req.Trailer = http.Header{"Apogee-Crc32": {"e2fc055d"}}
		res, err := newClient(t, startServer(t, s)).Do(req)
		if err != nil {
			t.Fatalf("PUT /recordings/furrow-a/YST/18856 failed: %v", err)
		}
		res.Body.Close()
		got, _ := s.Store.Get("furrow-a/YST/18856")
		if res.StatusCode != 201 || string(got) != "first piece second piece" {
			t.Errorf("PUT /recordings/furrow-a/YST/18856 = %d, keeping %q; want 201, keeping %q", res.StatusCode, got, "first piece second piece")
		}
	})
}

// chunked writes data as a chunked body in chunks of the given sizes, then the zero chunk
// and an empty trailer.
func chunked(data []byte, sizes ...int) string {
	var b strings.Builder
	for _, n := range sizes {
		fmt.Fprintf(&b, "%x\r\n%s\r\n", n, data[:n])
		data = data[n:]
	}
	b.WriteString("0\r\n\r\n")
	return b.String()
}

// pieces is a body that arrives one piece per read, with no length known in advance, so
// Go's client sends it chunked.
type pieces []string

// Read gives the next piece, then io.EOF.
func (p *pieces) Read(b []byte) (int, error) {
	if len(*p) == 0 {
		return 0, io.EOF
	}
	n := copy(b, (*p)[0])
	*p = (*p)[1:]
	return n, nil
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

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestBodyLength(t *testing.T) {
	tests := []struct {
		name, method string
		fields       []string
		want         int64
		status       int
		msg          string
	}{
		{"chunked", "PUT", []string{"Transfer-Encoding: chunked"}, -1, 0, ""},
		{"Chunked in capitals", "PUT", []string{"Transfer-Encoding: Chunked"}, -1, 0, ""},
		{"gzip alone", "PUT", []string{"Transfer-Encoding: gzip"}, 0, 400, `bad transfer-encoding "gzip"`},
		{"chunked, then gzip", "PUT", []string{"Transfer-Encoding: chunked, gzip"}, 0, 400, `bad transfer-encoding "chunked, gzip"`},
		{"gzip, then chunked", "PUT", []string{"Transfer-Encoding: gzip, chunked"}, 0, 501, `transfer-encoding "gzip, chunked" not implemented`},
		{"chunked twice", "PUT", []string{"Transfer-Encoding: chunked", "Transfer-Encoding: chunked"}, 0, 501, `transfer-encoding "chunked, chunked" not implemented`},
		{"both, the length first", "PUT", []string{"Content-Length: 46", "Transfer-Encoding: chunked"}, 0, 400, "both content-length and transfer-encoding"},
		{"both, the coding first", "PUT", []string{"Transfer-Encoding: chunked", "Content-Length: 46"}, 0, 400, "both content-length and transfer-encoding"},
		{"both, the coding unknown", "PUT", []string{"Content-Length: 46", "Transfer-Encoding: gzip"}, 0, 400, "both content-length and transfer-encoding"},
		{"GET, chunked", "GET", []string{"Transfer-Encoding: chunked"}, 0, 400, "GET has no body"},
		{"HEAD, chunked", "HEAD", []string{"Transfer-Encoding: chunked"}, 0, 400, "HEAD has no body"},
		{"a length alone", "PUT", []string{"Content-Length: 4096"}, 4096, 0, ""},
		{"GET with neither", "GET", nil, 0, 0, ""},
		{"PUT with neither", "PUT", nil, 0, 411, "length required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := header(t, tt.fields...)
			call := fmt.Sprintf("bodyLength(%s, %q)", tt.method, tt.fields)
			got, err := bodyLength(tt.method, h, 65536)
			if tt.status == 0 {
				if got != tt.want || err != nil {
					t.Errorf("%s = %d, %v, want %d, nil", call, got, err, tt.want)
				}
				return
			}
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != tt.status || e.Msg != tt.msg {
				t.Errorf("%s = %d, %v, want %d %s", call, got, err, tt.status, tt.msg)
			}
		})
	}
	t.Run("Mirela's two requests in one", func(t *testing.T) {
		hidden := "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n"
		body := "0\r\n\r\n" + hidden
		raw := "PUT /recordings/furrow-a/YST/18644 HTTP/1.1\r\nHost: 127.0.0.1\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
			"Transfer-Encoding: chunked\r\n\r\n" + body
		s := &Server{Store: NewStore(), Now: june7, Limits: Limits{MaxBody: 65536, Read: 50 * time.Millisecond}}
		res, rest := roundTrip(t, startServer(t, s), raw)
		if res.StatusCode != 400 || !res.Close {
			t.Errorf("PUT with both lengths = %d (Connection: close %v), want 400 (Connection: close true)", res.StatusCode, res.Close)
		}
		if rest != "" {
			t.Errorf("after its answer, the server sent %q, want nothing more", rest)
		}
		if lines := s.Store.List(); len(lines) != 0 {
			t.Errorf("the store keeps %q, want nothing", lines)
		}
	})
}

// header reads the given field lines as a request's header, as the server would.
func header(t *testing.T, fields ...string) Header {
	t.Helper()
	text := strings.Join(fields, "\r\n")
	if text != "" {
		text += "\r\n"
	}
	h, err := readHeader(bufio.NewReader(strings.NewReader(text + "\r\n")))
	if err != nil {
		t.Fatalf("readHeader(%q) = %v; the test's fields should read", fields, err)
	}
	return h
}

// june7 is the station server's clock in these tests.
func june7() time.Time { return time.Date(2028, 6, 7, 14, 5, 0, 0, time.UTC) }

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
	io.ReadAll(res.Body)
	res.Body.Close()
	rest, err := io.ReadAll(br)
	if err != nil {
		t.Errorf("after the answer, the connection gave %v, want the server to hang up", err)
	}
	return res, string(rest)
}

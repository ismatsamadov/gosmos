package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

// goGET is the head Go's own client sends for GET http://127.0.0.1/recordings.
const goGET = "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n"

func TestReadRequest(t *testing.T) {
	upload := "PUT /recordings/furrow-a/YST/18218 HTTP/1.1\r\nHost: 127.0.0.1\r\n"
	tests := []struct {
		name   string
		head   string
		want   *Request // nil when the head is refused
		status int
		msg    string
	}{
		{"Go's own GET", goGET, &Request{Method: "GET", Target: "/recordings", Header: Header{
			"host": {"127.0.0.1"}, "user-agent": {"Go-http-client/1.1"}, "accept-encoding": {"gzip"},
		}}, 0, ""},
		{"an upload with its length", upload + "Content-Length: 4096\r\n\r\n", &Request{Method: "PUT", Target: "/recordings/furrow-a/YST/18218", Header: Header{
			"host": {"127.0.0.1"}, "content-length": {"4096"},
		}, Length: 4096}, 0, ""},
		{"an 1100-byte request line", "GET /" + strings.Repeat("a", 1095) + " HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n", nil, 414, "request line over 1024 bytes"},
		{"no Host", "GET /recordings HTTP/1.1\r\nUser-Agent: Go-http-client/1.1\r\n\r\n", nil, 400, "missing host"},
		{"two Host lines", "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\nHost: 127.0.0.1\r\n\r\n", nil, 400, "more than one host"},
		{"PATCH", "PATCH /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n", nil, 501, `method "PATCH" not implemented`},
		{"get in lower case", "get /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n", nil, 501, `method "get" not implemented`},
		{"an upload with no length", upload + "\r\n", nil, 411, "length required"},
		{"HTTP/1.0", "GET /recordings HTTP/1.0\r\nHost: 127.0.0.1\r\n\r\n", nil, 505, `version "HTTP/1.0" not supported`},
		{"33 fields", "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n" + fields(32) + "\r\n", nil, 431, "more than 32 header fields"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readRequest(bufio.NewReaderSize(strings.NewReader(tt.head), 1024), 65536)
			if tt.want != nil {
				if err != nil || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("readRequest(%s) = %+v, %v, want %+v, nil", tt.name, got, err, tt.want)
				}
				return
			}
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != tt.status || e.Msg != tt.msg {
				t.Errorf("readRequest(%s) = %+v, %v, want the RequestError %d %s", tt.name, got, err, tt.status, tt.msg)
			}
		})
	}
	t.Run("nothing at all", func(t *testing.T) {
		if _, err := readRequest(bufio.NewReaderSize(strings.NewReader(""), 1024), 65536); !errors.Is(err, io.EOF) {
			t.Errorf("readRequest(no bytes) = %v, want io.EOF: a client that says nothing gets no answer", err)
		}
	})
}

func TestReadRequestThroughServer(t *testing.T) {
	addr := startServer(t)
	t.Run("Go's own GET", func(t *testing.T) {
		client := &http.Client{Timeout: time.Second}
		defer client.CloseIdleConnections()
		res, err := client.Get("http://" + addr + "/recordings")
		if err != nil {
			t.Fatalf("GET /recordings: %v", err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		want := "cirrine-1/YST/13094 2048 9f5edd58\nfurrow-a/YST/18218 4096 a2912082\n"
		if res.StatusCode != 200 || string(body) != want || err != nil {
			t.Errorf("GET /recordings = %d %q, %v, want 200 %q", res.StatusCode, body, err, want)
		}
	})
	refusals := []struct {
		name, send string
		status     int
		msg        string
	}{
		{"PATCH /recordings", "PATCH /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n", 501, `method "PATCH" not implemented`},
		{"no Host", "GET /recordings HTTP/1.1\r\nUser-Agent: Go-http-client/1.1\r\n\r\n", 400, "missing host"},
		{"two Host lines", "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\nHost: 127.0.0.1\r\n\r\n", 400, "more than one host"},
		// The first 1024 bytes of an 1100-byte request line, and nothing after them: the
		// server refuses before the rest is sent.
		{"an 1100-byte request line", "GET /" + strings.Repeat("a", 1019), 414, "request line over 1024 bytes"},
		{"33 fields", "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n" + fields(32), 431, "more than 32 header fields"},
		{"half a head, then silence", "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n", 408, "request timeout"},
	}
	for _, r := range refusals {
		t.Run(r.name, func(t *testing.T) {
			conn, err := net.DialTimeout("tcp", addr, time.Second)
			if err != nil {
				t.Fatalf("dial: %v", err)
			}
			defer conn.Close()
			conn.SetDeadline(time.Now().Add(time.Second))
			if _, err := io.WriteString(conn, r.send); err != nil {
				t.Fatalf("send: %v", err)
			}
			br := bufio.NewReader(conn)
			res, err := http.ReadResponse(br, nil)
			if err != nil {
				t.Fatalf("%s: reading the answer: %v", r.name, err)
			}
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode != r.status || string(body) != r.msg+"\n" {
				t.Errorf("%s = %d %q, want %d %q", r.name, res.StatusCode, body, r.status, r.msg+"\n")
			}
			if !res.Close {
				t.Errorf("%s: the answer has no Connection: close, want it on every refusal", r.name)
			}
			if _, err := br.ReadByte(); err != io.EOF {
				t.Errorf("%s: after the answer, a read = %v, want io.EOF: the server hangs up", r.name, err)
			}
		})
	}
}

// fields returns n header lines, X-Field-1 to X-Field-n, each ending in CR LF.
func fields(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "X-Field-%d: %d\r\n", i, i)
	}
	return b.String()
}

// testRecording returns an n-byte test recording: byte i is i % 256, so it holds zero
// bytes, CRs and LFs.
func testRecording(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 256)
	}
	return b
}

// startServer runs the station server on a loopback port the system picks, holding two
// recordings and giving a head 50 ms to arrive. When the test ends it closes the listener
// and waits for the server to stop.
func startServer(t *testing.T) string {
	t.Helper()
	store := NewStore()
	store.Put("cirrine-1/YST/13094", testRecording(2048))
	store.Put("furrow-a/YST/18218", testRecording(4096))
	srv := &Server{
		Store:  store,
		Limits: Limits{MaxBody: 65536, Head: 50 * time.Millisecond, Write: time.Second},
		Now:    func() time.Time { return time.Date(2028, 5, 26, 14, 0, 0, 0, time.UTC) },
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	t.Cleanup(func() {
		ln.Close()
		select {
		case err := <-served:
			if err != nil {
				t.Errorf("Serve = %v, want nil", err)
			}
		case <-time.After(time.Second):
			t.Error("the server did not stop within a second of its listener closing")
		}
	})
	return ln.Addr().String()
}

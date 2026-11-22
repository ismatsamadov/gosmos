package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestReadChunked(t *testing.T) {
	cuts := []struct{ name, stream string }{
		{"nothing after the head", ""},
		{"after a whole chunk", "c\r\nfirst piece \r\n"},
		{"inside a size line", "c\r\nfirst piece \r\nc"},
		{"inside a chunk", "c\r\nfirst piece \r\nc\r\nsecond"},
		{"after the zero chunk", "c\r\nfirst piece \r\nc\r\nsecond piece\r\n0\r\n"},
		{"inside the trailer", "c\r\nfirst piece \r\nc\r\nsecond piece\r\n0\r\nApogee-Crc32: e2fc"},
		{"before the last blank line", "c\r\nfirst piece \r\nc\r\nsecond piece\r\n0\r\nApogee-Crc32: e2fc055d\r\n"},
	}
	for _, tt := range cuts {
		t.Run("a stream ending "+tt.name, func(t *testing.T) {
			body, _, err := readChunked(bufio.NewReaderSize(strings.NewReader(tt.stream), 1024), 65536)
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != 400 || e.Msg != "recording cut off" {
				t.Errorf("readChunked(a stream ending %s) = %q, %v, want 400 recording cut off", tt.name, body, err)
			}
		})
	}
	t.Run("a whole stream", func(t *testing.T) {
		stream := "c\r\nfirst piece \r\nc\r\nsecond piece\r\n0\r\nApogee-Crc32: e2fc055d\r\n\r\n"
		body, trailer, err := readChunked(bufio.NewReaderSize(strings.NewReader(stream), 1024), 65536)
		if string(body) != "first piece second piece" || trailer.Get("apogee-crc32") != "e2fc055d" || err != nil {
			t.Errorf("readChunked(a whole stream) = %q, apogee-crc32 %q, %v, want %q, %q, nil",
				body, trailer.Get("apogee-crc32"), err, "first piece second piece", "e2fc055d")
		}
	})
	t.Run("Go's client through the station server", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june6, Limits: Limits{MaxBody: 65536, Read: 50 * time.Millisecond}}
		addr, stop := startServer(t, s)
		client := newClient(t, addr)

		whole, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/pipit-2/YST/16592", &pieces{"first piece ", "second piece"})
		if err != nil {
			t.Fatal(err)
		}
		res, err := client.Do(whole)
		if err != nil {
			t.Fatalf("a whole upload of pipit-2/YST/16592 failed: %v", err)
		}
		res.Body.Close()
		if res.StatusCode != 201 {
			t.Errorf("PUT /recordings/pipit-2/YST/16592 = %d, want 201", res.StatusCode)
		}

		cut, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/furrow-a/YST/18644", &pieces{"first piece ", "!recorder restarted"})
		if err != nil {
			t.Fatal(err)
		}
		if res, err := client.Do(cut); err == nil {
			res.Body.Close()
			t.Fatalf("an upload whose body failed got %d; the test's body should fail", res.StatusCode)
		}

		res, err = client.Get("http://127.0.0.1/recordings")
		if err != nil {
			t.Fatalf("GET /recordings failed: %v", err)
		}
		list, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if want := "pipit-2/YST/16592 24 e2fc055d\n"; string(list) != want {
			t.Errorf("GET /recordings after the cut upload = %q, want %q", list, want)
		}
		stop() // every connection has been served when it returns
		if got, ok := s.Store.Get("furrow-a/YST/18644"); ok {
			t.Errorf("the store keeps %q from the cut upload, want nothing", got)
		}
	})
}

// pieces is a body that arrives one piece per read, with no length known in advance, so
// Go's client sends it chunked. A piece starting with ! is an error instead, as when the
// recorder restarts.
type pieces []string

// Read gives the next piece, or its error, then io.EOF.
func (p *pieces) Read(b []byte) (int, error) {
	if len(*p) == 0 {
		return 0, io.EOF
	}
	piece := (*p)[0]
	*p = (*p)[1:]
	if msg, ok := strings.CutPrefix(piece, "!"); ok {
		return 0, errors.New(msg)
	}
	return copy(b, piece), nil
}

// june6 is the station server's clock in these tests.
func june6() time.Time { return time.Date(2028, 6, 6, 8, 6, 0, 0, time.UTC) }

// startServer serves s on a loopback port the system picks. It returns the address and a
// stop function, which closes the listener and waits until every connection has been
// served; the test's end calls it too.
func startServer(t *testing.T, s *Server) (string, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	served := make(chan error, 1)
	go func() { served <- s.Serve(ln) }()
	var once sync.Once
	stop := func() {
		once.Do(func() {
			ln.Close()
			select {
			case <-served:
			case <-time.After(time.Second):
				t.Error("the station server was still serving a second after its listener closed")
			}
		})
	}
	t.Cleanup(stop)
	return ln.Addr().String(), stop
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

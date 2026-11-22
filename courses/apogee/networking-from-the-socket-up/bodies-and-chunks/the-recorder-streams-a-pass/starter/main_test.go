package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestPut(t *testing.T) {
	s := &Server{Store: NewStore(), Now: june9, Limits: Limits{MaxBody: 65536, Read: 50 * time.Millisecond}}
	client := newClient(t, startServer(t, s))
	furrowA := testRecording(4096)
	uploads := []struct {
		name, id string
		body     []byte
		chunked  bool   // sent chunked, Apogee-Crc32 in the trailer; else with a length, in the header
		crc      string // the Apogee-Crc32 sent; empty for none
		status   int
		answer   string
		crcBack  string // the Apogee-Crc32 a 201 answers with
	}{
		{"Furrow-A's pass, chunked", "furrow-a/YST/18856", furrowA, true, "a2912082", 201, "stored furrow-a/YST/18856 4096 bytes\n", "a2912082"},
		{"the same pass again", "furrow-a/YST/18856", furrowA, true, "a2912082", 409, "furrow-a/YST/18856 already stored\n", ""},
		{"Cirrine-1's pass, with a length and no CRC", "cirrine-1/YST/14140", testRecording(2048), false, "", 201, "stored cirrine-1/YST/14140 2048 bytes\n", "9f5edd58"},
		{"a pass with the wrong CRC", "furrow-b/YST/18688", furrowA, true, "a2912083", 400, "crc mismatch: sent a2912083, received a2912082\n", ""},
		{"an orbit number of four digits", "furrow-a/YST/1885", furrowA, true, "a2912082", 404, "no such path\n", ""},
		{"a station code in small letters", "furrow-a/yst/18856", furrowA, true, "a2912082", 404, "no such path\n", ""},
	}
	for _, u := range uploads {
		t.Run(u.name, func(t *testing.T) {
			var body io.Reader = bytes.NewReader(u.body)
			if u.chunked {
				body = io.MultiReader(body) // no length known in advance: Go's client sends it chunked
			}
			req, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/"+u.id, body)
			if err != nil {
				t.Fatal(err)
			}
			switch {
			case u.crc == "":
			case u.chunked:
				req.Trailer = http.Header{"Apogee-Crc32": {u.crc}}
			default:
				req.Header.Set("Apogee-Crc32", u.crc)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("PUT /recordings/%s failed: %v", u.id, err)
			}
			got, _ := io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode != u.status || string(got) != u.answer {
				t.Errorf("PUT /recordings/%s = %d %q, want %d %q", u.id, res.StatusCode, got, u.status, u.answer)
			}
			if u.status != 201 {
				return
			}
			location, crc := res.Header.Get("Location"), res.Header.Get("Apogee-Crc32")
			if location != "/recordings/"+u.id || crc != u.crcBack {
				t.Errorf("PUT /recordings/%s answered Location %q, Apogee-Crc32 %q; want %q, %q",
					u.id, location, crc, "/recordings/"+u.id, u.crcBack)
			}
		})
	}
	t.Run("a cut upload", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/furrow-a/YST/18644", &cut{})
		if err != nil {
			t.Fatal(err)
		}
		if res, err := client.Do(req); err == nil {
			res.Body.Close()
			t.Fatalf("a cut upload got %d; the test's body should fail", res.StatusCode)
		}
	})
	t.Run("the list", func(t *testing.T) {
		res, err := client.Get("http://127.0.0.1/recordings")
		if err != nil {
			t.Fatalf("GET /recordings failed: %v", err)
		}
		got, _ := io.ReadAll(res.Body)
		res.Body.Close()
		want := "cirrine-1/YST/14140 2048 9f5edd58\nfurrow-a/YST/18856 4096 a2912082\n"
		if string(got) != want {
			t.Errorf("GET /recordings = %q, want %q", got, want)
		}
	})
}

// cut is a body that sends one piece, then fails, as when the recorder restarts mid-pass.
type cut struct{ sent bool }

// Read gives one piece, then an error.
func (c *cut) Read(b []byte) (int, error) {
	if c.sent {
		return 0, errors.New("recorder restarted")
	}
	c.sent = true
	return copy(b, strings.Repeat("furrow ", 100)), nil
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

// june9 is the station server's clock in these tests.
func june9() time.Time { return time.Date(2028, 6, 9, 10, 51, 0, 0, time.UTC) }

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

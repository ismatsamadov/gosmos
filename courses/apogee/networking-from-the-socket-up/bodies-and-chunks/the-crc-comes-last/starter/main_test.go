package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestCheckCRC(t *testing.T) {
	pieces := []byte("first piece second piece")
	recording := testRecording(4096)
	tests := []struct {
		name            string
		body            []byte
		header, trailer Header
		want            uint32
		status          int
		msg             string
	}{
		{"the trailer's CRC", pieces, Header{"trailer": {"Apogee-Crc32"}}, Header{"apogee-crc32": {"e2fc055d"}}, 0xe2fc055d, 0, ""},
		{"a trailer one digit off", pieces, nil, Header{"apogee-crc32": {"e2fc055e"}}, 0, 400, "crc mismatch: sent e2fc055e, received e2fc055d"},
		{"the header's CRC", recording, Header{"apogee-crc32": {"a2912082"}}, nil, 0xa2912082, 0, ""},
		{"a header one digit off", recording, Header{"apogee-crc32": {"a2912083"}}, nil, 0, 400, "crc mismatch: sent a2912083, received a2912082"},
		{"the trailer over the header", pieces, Header{"apogee-crc32": {"00000000"}}, Header{"apogee-crc32": {"e2fc055d"}}, 0xe2fc055d, 0, ""},
		{"no CRC at all", recording, nil, nil, 0xa2912082, 0, ""},
		{"capital hex digits", recording, Header{"apogee-crc32": {"A2912082"}}, nil, 0xa2912082, 0, ""},
		{"not hex", recording, Header{"apogee-crc32": {"xyz"}}, nil, 0, 400, `bad apogee-crc32 "xyz"`},
		{"seven digits", pieces, nil, Header{"apogee-crc32": {"e2fc055"}}, 0, 400, `bad apogee-crc32 "e2fc055"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := fmt.Sprintf("checkCRC(%s)", tt.name)
			got, err := checkCRC(tt.body, tt.header, tt.trailer)
			if tt.status == 0 {
				if got != tt.want || err != nil {
					t.Errorf("%s = %08x, %v, want %08x, nil", call, got, err, tt.want)
				}
				return
			}
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != tt.status || e.Msg != tt.msg {
				t.Errorf("%s = %08x, %v, want %d %s", call, got, err, tt.status, tt.msg)
			}
		})
	}
}

func TestCRCThroughTheServer(t *testing.T) {
	s := &Server{Store: NewStore(), Now: june8, Limits: Limits{MaxBody: 65536, Read: 50 * time.Millisecond}}
	client := newClient(t, startServer(t, s))
	uploads := []struct {
		name, id string
		body     io.Reader
		header   string // Apogee-Crc32 in the head
		trailer  string // Apogee-Crc32 after the last chunk
		status   int
		kept     bool
	}{
		{"chunked, its CRC in the trailer", "furrow-a/YST/18674", &pieces{"first piece ", "second piece"}, "", "e2fc055d", 201, true},
		{"chunked, the trailer one digit off", "furrow-b/YST/18672", &pieces{"first piece ", "second piece"}, "", "e2fc055e", 400, false},
		{"Content-Length, its CRC in the header", "pipit-2/YST/16592", bytes.NewReader(testRecording(1024)), "b70b4c26", "", 201, true},
		{"Content-Length, the header one digit off", "pipit-2/YST/16637", bytes.NewReader(testRecording(1024)), "b70b4c27", "", 400, false},
	}
	for _, u := range uploads {
		t.Run(u.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/"+u.id, u.body)
			if err != nil {
				t.Fatal(err)
			}
			if u.header != "" {
				req.Header.Set("Apogee-Crc32", u.header)
			}
			if u.trailer != "" {
				req.Trailer = http.Header{"Apogee-Crc32": {u.trailer}}
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("PUT /recordings/%s failed: %v", u.id, err)
			}
			msg, _ := io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode != u.status {
				t.Errorf("PUT /recordings/%s = %d %q, want %d", u.id, res.StatusCode, msg, u.status)
			}
			if _, kept := s.Store.Get(u.id); kept != u.kept {
				t.Errorf("after PUT /recordings/%s, the store keeps it: %v, want %v", u.id, kept, u.kept)
			}
		})
	}
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

// june8 is the station server's clock in these tests.
func june8() time.Time { return time.Date(2028, 6, 8, 10, 52, 0, 0, time.UTC) }

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

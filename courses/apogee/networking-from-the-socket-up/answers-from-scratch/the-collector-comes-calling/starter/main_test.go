package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestRoute(t *testing.T) {
	s := &Server{Store: NewStore(), Now: june20, Limits: Limits{MaxBody: 65536}}
	s.Store.Put("furrow-a/YST/18856", testRecording(4096))
	client, _ := newClient(t, startServer(t, s))
	tests := []struct {
		method, path string
		status       int
		allow        string
		body         string
	}{
		{"GET", "/recordings", 200, "", "furrow-a/YST/18856 4096 a2912082\n"},
		{"HEAD", "/recordings", 200, "", ""},
		{"PUT", "/recordings", 405, "GET, HEAD", "method not allowed\n"},
		{"GET", "/recordings/furrow-a/YST/18856", 200, "", string(testRecording(4096))},
		{"GET", "/recordings/furrow-a/YST/18857", 404, "", "no recording furrow-a/YST/18857\n"},
		{"GET", "/recordings/furrow-a/YST/1885", 404, "", "no such path\n"},
		{"HEAD", "/feed", 405, "GET", ""},
		{"PUT", "/feed", 405, "GET", "method not allowed\n"},
		{"GET", "/recordings?sat=furrow-a", 404, "", "no such path\n"},
		{"GET", "/weather/YST", 404, "", "no such path\n"},
		{"DELETE", "/recordings/furrow-a/YST/18856", 501, "", "method \"DELETE\" not implemented\n"},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var body io.Reader
			if tt.method == "PUT" {
				body = strings.NewReader("")
			}
			req, err := http.NewRequest(tt.method, "http://127.0.0.1"+tt.path, body)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("%s %s failed: %v", tt.method, tt.path, err)
			}
			got, _ := io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode != tt.status || res.Header.Get("Allow") != tt.allow || string(got) != tt.body {
				t.Errorf("%s %s = %d, Allow %q, %s; want %d, Allow %q, %s",
					tt.method, tt.path, res.StatusCode, res.Header.Get("Allow"), clip(got), tt.status, tt.allow, clip([]byte(tt.body)))
			}
		})
	}
}

func TestTheCollectorComesCalling(t *testing.T) {
	s := &Server{Store: NewStore(), Now: june20, Limits: Limits{MaxBody: 65536, Head: time.Second, Read: time.Second}}
	addr, served := startServerFor(t, s)
	collector, dials := newClient(t, addr)
	recorder, _ := newClient(t, addr)

	feed, err := collector.Get("http://127.0.0.1/feed")
	if err != nil {
		t.Fatalf("the collector's GET /feed failed: %v", err)
	}
	defer feed.Body.Close()
	if feed.StatusCode != 200 || feed.ContentLength != -1 {
		t.Fatalf("GET /feed = %d with ContentLength %d, want 200, chunked (-1)", feed.StatusCode, feed.ContentLength)
	}
	lines := bufio.NewReader(feed.Body)

	furrowA := testRecording(4096)
	req, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/furrow-a/YST/18856", io.MultiReader(bytes.NewReader(furrowA)))
	if err != nil {
		t.Fatal(err)
	}
	req.Trailer = http.Header{"Apogee-Crc32": {"a2912082"}}
	res, err := recorder.Do(req)
	if err != nil {
		t.Fatalf("the recorder's upload of furrow-a/YST/18856 failed: %v", err)
	}
	res.Body.Close()
	if res.StatusCode != 201 {
		t.Fatalf("PUT /recordings/furrow-a/YST/18856 = %d, want 201", res.StatusCode)
	}
	if line, err := lines.ReadString('\n'); line != "furrow-a/YST/18856 4096 a2912082\n" || err != nil {
		t.Fatalf("the feed gave %q, %v, want %q", line, err, "furrow-a/YST/18856 4096 a2912082\n")
	}

	for _, method := range []string{"HEAD", "GET"} {
		req, err := http.NewRequest(method, "http://127.0.0.1/recordings/furrow-a/YST/18856", nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := collector.Do(req)
		if err != nil {
			t.Fatalf("the collector's %s of furrow-a/YST/18856 failed: %v", method, err)
		}
		got, _ := io.ReadAll(res.Body)
		res.Body.Close()
		crc := fmt.Sprintf("%08x", crc32.ChecksumIEEE(got))
		switch {
		case method == "HEAD" && (res.StatusCode != 200 || res.ContentLength != 4096):
			t.Errorf("HEAD of furrow-a/YST/18856 = %d with ContentLength %d, want 200 with 4096", res.StatusCode, res.ContentLength)
		case method == "GET" && (res.StatusCode != 200 || len(got) != 4096 || crc != res.Header.Get("Apogee-Crc32")):
			t.Errorf("GET of furrow-a/YST/18856 = %d, %d bytes of crc %s, Apogee-Crc32 %q; want 200, 4096 bytes whose crc it names",
				res.StatusCode, len(got), crc, res.Header.Get("Apogee-Crc32"))
		}
	}
	if n := dials.Load(); n != 2 {
		t.Errorf("the collector made %d connections, want 2: the feed's, and one for HEAD and GET", n)
	}

	cut := dial(t, addr)
	io.WriteString(cut, "PUT /recordings/cirrine-1/YST/14302 HTTP/1.1\r\nHost: 127.0.0.1\r\nTransfer-Encoding: chunked\r\n\r\n400\r\n"+string(testRecording(1024))+"\r\n")
	cut.(*net.TCPConn).CloseWrite() // the recorder restarts mid-stream
	if res, err := http.ReadResponse(bufio.NewReader(cut), nil); err != nil || res.StatusCode != 400 {
		t.Fatalf("Cirrine-1's cut upload = %v, %v, want 400", res, err)
	}
	again, err := http.NewRequest("PUT", "http://127.0.0.1/recordings/cirrine-1/YST/14302", bytes.NewReader(testRecording(2048)))
	if err != nil {
		t.Fatal(err)
	}
	if res, err := recorder.Do(again); err != nil || res.StatusCode != 201 {
		t.Fatalf("Cirrine-1's upload sent again whole = %v, %v, want 201", res, err)
	} else {
		res.Body.Close()
	}
	if line, err := lines.ReadString('\n'); line != "cirrine-1/YST/14302 2048 9f5edd58\n" || err != nil {
		t.Fatalf("after the cut upload and the whole one, the feed gave %q, %v, want only the whole one, %q", line, err, "cirrine-1/YST/14302 2048 9f5edd58\n")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown = %v, want nil", err)
	}
	if rest, err := io.ReadAll(lines); len(rest) != 0 || err != nil {
		t.Errorf("after Shutdown the feed gave %q, %v, want its last chunk: nothing more and no error", rest, err)
	}
	select {
	case err := <-served:
		if err != nil {
			t.Errorf("Serve = %v after Shutdown, want nil", err)
		}
	case <-time.After(time.Second):
		t.Error("Serve has not returned a second after Shutdown")
	}
}

// clip shortens a long body for a failure message.
func clip(b []byte) string {
	if len(b) <= 40 {
		return fmt.Sprintf("%q", b)
	}
	return fmt.Sprintf("%d bytes (crc %08x)", len(b), crc32.ChecksumIEEE(b))
}

// dial connects to the server at addr, with a second's deadline, closed when the test ends.
func dial(t *testing.T, addr string) net.Conn {
	t.Helper()
	c, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { c.Close() })
	c.SetDeadline(time.Now().Add(time.Second))
	return c
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

// june20 is the station server's clock in these tests: Furrow-A's LOS.
func june20() time.Time { return time.Date(2028, 6, 20, 6, 40, 35, 0, time.UTC) }

// startServer serves s on a loopback port the system picks and returns its address. When
// the test ends, it closes the listener and the store, and waits until Serve returns.
func startServer(t *testing.T, s *Server) string {
	t.Helper()
	addr, _ := startServerFor(t, s)
	return addr
}

// startServerFor is startServer that also returns the channel Serve's result arrives on.
func startServerFor(t *testing.T, s *Server) (string, <-chan error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	served, result := make(chan error, 1), make(chan error, 1)
	go func() {
		err := s.Serve(ln)
		served <- err
		result <- err
	}()
	t.Cleanup(func() {
		ln.Close()
		s.Store.Close()
		select {
		case <-served:
		case <-time.After(time.Second):
			t.Error("the station server was still serving a second after its listener closed")
		}
	})
	return ln.Addr().String(), result
}

// newClient returns an HTTP client whose every connection goes to addr, so its requests can
// name http://127.0.0.1/, and the count of connections it has dialed.
func newClient(t *testing.T, addr string) (*http.Client, *atomic.Int32) {
	dials := new(atomic.Int32)
	tr := &http.Transport{DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		dials.Add(1)
		var d net.Dialer
		return d.DialContext(ctx, "tcp", addr)
	}}
	t.Cleanup(tr.CloseIdleConnections)
	return &http.Client{Timeout: time.Second, Transport: tr}, dials
}

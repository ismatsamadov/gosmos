package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

// uploadHead and uploadEnd are the recorder's chunked upload of Furrow-A's pass, before and
// after its one chunk of 4,096 bytes.
const (
	uploadHead = "PUT /recordings/furrow-a/YST/18856 HTTP/1.1\r\nHost: 127.0.0.1\r\nTransfer-Encoding: chunked\r\n\r\n1000\r\n"
	uploadEnd  = "\r\n0\r\nApogee-Crc32: a2912082\r\n\r\n"
)

func TestShutdown(t *testing.T) {
	t.Run("during Tuesday's passes", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536}}
		addr, served := startServer(t, s)

		feed, err := newClient(t, addr).Get("http://127.0.0.1/feed")
		if err != nil {
			t.Fatalf("GET /feed failed: %v", err)
		}
		defer feed.Body.Close()
		idle := dial(t, addr)
		send(t, idle, "GET /recordings HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
		if res := readAnswer(t, idle); res.StatusCode != 200 {
			t.Fatalf("GET /recordings = %d, want 200", res.StatusCode)
		}
		upload := dial(t, addr)
		send(t, upload, uploadHead+string(testRecording(4096)[:2048]))
		waitFor(t, "the feed and the upload to be busy and the other connection to wait", func() bool {
			busy, waiting := conns(s)
			return busy == 2 && waiting == 1
		})

		stopped := make(chan error, 1)
		go func() { stopped <- s.Shutdown(context.Background()) }()

		if err := hungUp(idle, 200*time.Millisecond); err != nil {
			t.Errorf("the connection waiting between requests: %v, want it closed at once", err)
		}
		if rest, err := endOf(feed.Body, 200*time.Millisecond); rest != "" || err != nil {
			t.Errorf("the feed gave %q, %v, want its last chunk at once: nothing more and no error", rest, err)
		}
		select {
		case err := <-stopped:
			t.Fatalf("Shutdown = %v while an upload was in flight, want it to wait for the upload's answer", err)
		default:
		}

		send(t, upload, string(testRecording(4096)[2048:])+uploadEnd)
		if res := readAnswer(t, upload); res.StatusCode != 201 {
			t.Errorf("the upload in flight = %d, want 201", res.StatusCode)
		}
		if err := hungUp(upload, 200*time.Millisecond); err != nil {
			t.Errorf("after the upload's answer: %v, want its connection closed", err)
		}
		select {
		case err := <-stopped:
			if err != nil {
				t.Errorf("Shutdown = %v, want nil", err)
			}
		case <-time.After(time.Second):
			t.Fatal("Shutdown has not returned a second after the upload's answer")
		}
		if _, ok := s.Store.Get("furrow-a/YST/18856"); !ok {
			t.Error("the store doesn't keep furrow-a/YST/18856, want the upload kept")
		}
		if c, err := net.DialTimeout("tcp", addr, time.Second); err == nil {
			c.Close()
			t.Error("a new connection was accepted after Shutdown, want it refused")
		}
		select {
		case err := <-served:
			if err != nil {
				t.Errorf("Serve = %v after Shutdown, want nil", err)
			}
		case <-time.After(time.Second):
			t.Error("Serve has not returned a second after Shutdown")
		}
	})
	t.Run("an upload that never ends, 100ms", func(t *testing.T) {
		s := &Server{Store: NewStore(), Now: june19, Limits: Limits{MaxBody: 65536}}
		addr, served := startServer(t, s)
		upload := dial(t, addr)
		send(t, upload, uploadHead+"furrow")
		waitFor(t, "the upload to be busy", func() bool {
			busy, _ := conns(s)
			return busy == 1
		})
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if err := s.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Shutdown(100ms) with an upload that never ends = %v, want context.DeadlineExceeded", err)
		}
		if err := hungUp(upload, 200*time.Millisecond); err != nil {
			t.Errorf("the unfinished upload's connection: %v, want it closed", err)
		}
		select {
		case err := <-served:
			if err != nil {
				t.Errorf("Serve = %v after Shutdown, want nil", err)
			}
		case <-time.After(200 * time.Millisecond):
			t.Error("Serve has not returned after Shutdown")
		}
	})
}

// conns counts the server's open connections: those busy with a request, and those waiting
// between requests.
func conns(s *Server) (busy, waiting int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range s.conns {
		if b {
			busy++
		} else {
			waiting++
		}
	}
	return busy, waiting
}

// waitFor waits up to a second for cond to hold, and stops the test if it never does.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	for deadline := time.Now().Add(time.Second); !cond(); time.Sleep(time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatalf("waited a second for %s", what)
		}
	}
}

// rawConn is a test's own connection to the server and the reader over it.
type rawConn struct {
	net.Conn
	br *bufio.Reader
}

// dial connects to the server at addr, with a second's deadline, closed when the test ends.
func dial(t *testing.T, addr string) *rawConn {
	t.Helper()
	c, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { c.Close() })
	c.SetDeadline(time.Now().Add(time.Second))
	return &rawConn{Conn: c, br: bufio.NewReader(c)}
}

// send writes raw on c.
func send(t *testing.T, c *rawConn, raw string) {
	t.Helper()
	if _, err := io.WriteString(c, raw); err != nil {
		t.Fatalf("sending: %v", err)
	}
}

// readAnswer reads one answer from c, body and all.
func readAnswer(t *testing.T, c *rawConn) *http.Response {
	t.Helper()
	res, err := http.ReadResponse(c.br, nil)
	if err != nil {
		t.Fatalf("reading the answer: %v", err)
	}
	io.ReadAll(res.Body)
	res.Body.Close()
	return res
}

// hungUp waits up to d for the server to close c, and says what came instead if it didn't.
func hungUp(c *rawConn, d time.Duration) error {
	c.SetReadDeadline(time.Now().Add(d))
	if _, err := c.br.ReadByte(); err == nil {
		return errors.New("the server sent more bytes")
	} else if !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

// endOf reads what is left of body, waiting d at most, and returns it with the error the
// body ended with.
func endOf(body io.ReadCloser, d time.Duration) (string, error) {
	type result struct {
		rest []byte
		err  error
	}
	done := make(chan result, 1)
	go func() {
		rest, err := io.ReadAll(body)
		done <- result{rest, err}
	}()
	select {
	case r := <-done:
		return string(r.rest), r.err
	case <-time.After(d):
		body.Close()
		<-done
		return "", fmt.Errorf("no end after %v", d)
	}
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

// june19 is the station server's clock in these tests.
func june19() time.Time { return time.Date(2028, 6, 19, 22, 40, 0, 0, time.UTC) }

// startServer serves s on a loopback port the system picks and returns its address and the
// channel Serve's result arrives on. When the test ends, it closes the listener and the
// store, and waits a second at most for Serve to return.
func startServer(t *testing.T, s *Server) (string, <-chan error) {
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

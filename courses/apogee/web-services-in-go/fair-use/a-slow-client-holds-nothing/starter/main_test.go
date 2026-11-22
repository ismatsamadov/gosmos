package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

// stations answers GET /stations, like the booking service, with one station.
func stations() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"code":"YST","name":"Ystmark"}]` + "\n"))
	})
	return mux
}

// serve runs newServer(stations(), tm) on a 127.0.0.1 port the system picks, and returns its address.
func serve(t *testing.T, tm Timeouts) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(stations(), tm)
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close() })
	return ln.Addr().String()
}

// dial opens a connection to addr that gives up reading after half a second.
func dial(t *testing.T, addr string) net.Conn {
	t.Helper()
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
	return conn
}

// stillOpen reports whether reading from r waited out the connection's half-second deadline.
func stillOpen(r io.Reader) bool {
	_, err := r.Read(make([]byte, 1))
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func TestNewServer(t *testing.T) {
	h := http.NewServeMux()
	tm := Timeouts{ReadHeader: 5 * time.Second, Read: 15 * time.Second, Write: 30 * time.Second, Idle: 2 * time.Minute, Handler: 10 * time.Second}
	srv := newServer(h, tm)
	if srv.Handler != h {
		t.Errorf("newServer(h, t).Handler = %v, want h", srv.Handler)
	}
	for _, f := range []struct {
		name      string
		got, want time.Duration
	}{
		{"ReadHeaderTimeout", srv.ReadHeaderTimeout, tm.ReadHeader},
		{"ReadTimeout", srv.ReadTimeout, tm.Read},
		{"WriteTimeout", srv.WriteTimeout, tm.Write},
		{"IdleTimeout", srv.IdleTimeout, tm.Idle},
	} {
		if f.got != f.want {
			t.Errorf("newServer(h, t).%s = %v, want %v", f.name, f.got, f.want)
		}
	}
	if srv.MaxHeaderBytes != 16384 {
		t.Errorf("newServer(h, t).MaxHeaderBytes = %d, want 16384", srv.MaxHeaderBytes)
	}
}

func TestHalfARequestIsCut(t *testing.T) {
	t.Parallel()
	conn := dial(t, serve(t, Timeouts{ReadHeader: 50 * time.Millisecond}))
	conn.Write([]byte("GET /stations HTTP/1.1\r\nHost: 127.0.0.1\r\n")) // and never the blank line that ends the headers
	if stillOpen(conn) {
		t.Errorf("a connection that sent half its headers, at a 50ms ReadHeader timeout = open after half a second, want closed")
	}
}

func TestIdleConnectionIsClosed(t *testing.T) {
	t.Parallel()
	conn := dial(t, serve(t, Timeouts{Idle: 50 * time.Millisecond}))
	conn.Write([]byte("GET /stations HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n"))
	r := bufio.NewReader(conn)
	res, err := http.ReadResponse(r, nil)
	if err != nil {
		t.Fatalf("reading the answer to GET /stations: %v", err)
	}
	io.Copy(io.Discard, res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("GET /stations = %d, want 200", res.StatusCode)
	}
	if stillOpen(r) {
		t.Errorf("a kept-alive connection left idle, at a 50ms Idle timeout = open after half a second, want closed")
	}
}

func TestBigHeadersAreRefused(t *testing.T) {
	addr := serve(t, Timeouts{})
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+"/stations", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Note", strings.Repeat("a", 32<<10))
	res, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		t.Fatalf("GET /stations with a 32 KiB header: %v", err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusRequestHeaderFieldsTooLarge {
		t.Errorf("GET /stations with a 32 KiB header = %d, want 431", res.StatusCode)
	}
}

func TestOrdinaryRequestIsAnswered(t *testing.T) {
	second := time.Second
	addr := serve(t, Timeouts{ReadHeader: second, Read: second, Write: second, Idle: second})
	res, err := (&http.Client{Timeout: 2 * time.Second}).Get("http://" + addr + "/stations")
	if err != nil {
		t.Fatalf("GET /stations with every limit at 1s: %v", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil || res.StatusCode != http.StatusOK || !strings.Contains(string(body), "Ystmark") {
		t.Errorf("GET /stations with every limit at 1s = %d %q, %v, want 200 and the station list", res.StatusCode, body, err)
	}
}

package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

// listen returns a listener on 127.0.0.1 on a port the system picks, closed when the test ends.
func listen(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln
}

// get sends GET path to the listener's address and returns the status and the trimmed body, or the error.
func get(ln net.Listener, path string) (int, string, error) {
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Get("http://" + ln.Addr().String() + path)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return resp.StatusCode, strings.TrimSpace(string(body)), err
}

func TestDrainingBeforeTheStop(t *testing.T) {
	svc := &Service{Guide: fstest.MapFS{"index.html": {Data: []byte("<h1>Booking passes with Apogee</h1>\n")}}}
	mux := http.NewServeMux() // the two routes this lesson needs, as the service registers them
	mux.HandleFunc("GET /health", svc.health)
	mux.Handle("GET /guide/", http.StripPrefix("/guide/", http.FileServerFS(svc.Guide)))
	app := &App{Service: svc, Handler: mux}
	limits := Timeouts{Drain: 300 * time.Millisecond, Shutdown: time.Second}
	ln := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- run(ctx, ln, app, limits) }()

	if status, body, err := get(ln, "/health"); err != nil || status != http.StatusOK || body != `{"status":"ok"}` {
		t.Fatalf("GET /health before the stop = %d %s, %v, want 200 {\"status\":\"ok\"}", status, body, err)
	}
	stopped := time.Now()
	cancel()
	status, body, err := get(ln, "/health")
	for err == nil && status == http.StatusOK && time.Since(stopped) < time.Second {
		status, body, err = get(ln, "/health")
	}
	if err != nil || status != http.StatusServiceUnavailable || body != `{"status":"draining"}` {
		t.Fatalf("GET /health after the stop = %d %s, %v, want 503 {\"status\":\"draining\"} while the service drains", status, body, err)
	}
	if status, _, err := get(ln, "/guide/"); err != nil || status != http.StatusOK {
		t.Errorf("GET /guide/ while draining = %d, %v, want 200: every other route still works", status, err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("run = %v, want nil", err)
		}
		if took := time.Since(stopped); took < limits.Drain {
			t.Errorf("run returned %v after the stop, want it to drain for %v first", took, limits.Drain)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("run hadn't returned 2s after the stop")
	}
}

package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

// rehearsal is the service as the tests run it: a health check, and GET /bookings held until release is closed.
type rehearsal struct {
	app     *App
	entered chan struct{}
	release chan struct{}
}

// newRehearsal builds the rehearsal service.
func newRehearsal() *rehearsal {
	r := &rehearsal{entered: make(chan struct{}, 1), release: make(chan struct{})}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("GET /bookings", func(w http.ResponseWriter, req *http.Request) {
		r.entered <- struct{}{}
		select {
		case <-r.release:
			io.WriteString(w, `[{"id":"bk-4801"}]`)
		case <-req.Context().Done():
		}
	})
	r.app = &App{Handler: mux}
	return r
}

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

// get sends GET path to the listener's address and returns the status, or the error.
func get(ln net.Listener, path string) (int, error) {
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Get("http://" + ln.Addr().String() + path)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

// returned waits up to a second for run's result.
func returned(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(time.Second):
		t.Fatal("run hadn't returned a second later")
		return nil
	}
}

// stillRunning reports whether run has not returned yet, putting back a result it finds.
func stillRunning(done chan error) bool {
	select {
	case err := <-done:
		done <- err
		return false
	default:
		return true
	}
}

// limits are the rehearsal's timeouts: a second to shut down, nothing else limited.
var limits = Timeouts{Shutdown: time.Second}

func TestRunServesUntilCanceled(t *testing.T) {
	ln := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- run(ctx, ln, newRehearsal().app, limits) }()
	if status, err := get(ln, "/health"); err != nil || status != http.StatusOK {
		t.Fatalf("GET /health while running = %d, %v, want 200", status, err)
	}
	if !stillRunning(done) {
		t.Fatalf("run returned %v while ctx still lived, want it to keep serving", <-done)
	}
	cancel()
	if err := returned(t, done); err != nil {
		t.Errorf("run after ctx ended = %v, want nil", err)
	}
	if status, err := get(ln, "/health"); err == nil {
		t.Errorf("GET /health after run returned = %d, want the connection refused", status)
	}
}

func TestRunFinishesTheRequestInFlight(t *testing.T) {
	ln := listen(t)
	svc := newRehearsal()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- run(ctx, ln, svc.app, limits) }()
	status := make(chan int, 1)
	go func() {
		code, _ := get(ln, "/bookings")
		status <- code
	}()
	select {
	case <-svc.entered:
	case <-time.After(time.Second):
		t.Fatal("GET /bookings never reached the handler")
	}
	cancel()
	if !stillRunning(done) {
		t.Error("run returned while a request was in flight, want it to wait for the request")
	}
	close(svc.release)
	if got := <-status; got != http.StatusOK {
		t.Errorf("GET /bookings in flight at the stop = %d, want 200", got)
	}
	if err := returned(t, done); err != nil {
		t.Errorf("run after ctx ended = %v, want nil", err)
	}
}

func TestRunReturnsServeFailure(t *testing.T) {
	ln := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- run(ctx, ln, newRehearsal().app, limits) }()
	ln.Close() // the listener breaks underneath the server; ctx never ends
	if err := returned(t, done); !errors.Is(err, net.ErrClosed) {
		t.Errorf("run with its listener closed underneath = %v, want that failure, net.ErrClosed, without waiting for ctx", err)
	}
}

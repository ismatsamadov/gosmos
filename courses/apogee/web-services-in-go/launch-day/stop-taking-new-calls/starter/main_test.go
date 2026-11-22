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

// held is a slow GET /bookings: it holds every request until release is closed, after telling entered.
type held struct {
	entered chan struct{}
	release chan struct{}
}

// ServeHTTP waits for release, then answers with the bookings.
func (h *held) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.entered <- struct{}{}
	select {
	case <-h.release:
		io.WriteString(w, `[{"id":"bk-4801"}]`)
	case <-r.Context().Done():
	}
}

// serve starts a server for a new held handler on 127.0.0.1, on a port the system picks.
func serve(t *testing.T) (*http.Server, *held, string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	h := &held{entered: make(chan struct{}, 1), release: make(chan struct{})}
	srv := &http.Server{Handler: h}
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close() })
	return srv, h, ln.Addr().String()
}

// answer is what a client got back: a status and a body, or an error.
type answer struct {
	status int
	body   string
	err    error
}

// inFlight sends GET /bookings to addr in the background and passes on what came back.
func inFlight(addr string) <-chan answer {
	done := make(chan answer, 1)
	go func() {
		resp, err := (&http.Client{Timeout: 2 * time.Second}).Get("http://" + addr + "/bookings")
		if err != nil {
			done <- answer{err: err}
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		done <- answer{resp.StatusCode, string(body), err}
	}()
	return done
}

// wait waits up to a second for something on ch, and fails the test if nothing comes.
func wait[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(time.Second):
		t.Fatalf("waited a second for %s", what)
		var zero T
		return zero
	}
}

// refused reports whether addr refuses new connections within a second.
func refused(addr string) bool {
	for start := time.Now(); time.Since(start) < time.Second; time.Sleep(time.Millisecond) {
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err != nil {
			return true
		}
		conn.Close()
	}
	return false
}

func TestStopFinishesWhatIsInFlight(t *testing.T) {
	srv, h, addr := serve(t)
	res := inFlight(addr)
	wait(t, h.entered, "the request to reach the handler")
	stopped := make(chan error, 1)
	go func() { stopped <- stop(srv, time.Second) }()
	if !refused(addr) {
		t.Fatal("stop(srv, 1s): new connections were still taken a second later, want them refused at once")
	}
	select {
	case err := <-stopped:
		t.Errorf("stop(srv, 1s) = %v while a request was in flight, want it to wait for the request", err)
		stopped <- err
	default:
	}
	close(h.release)
	if got := wait(t, res, "the request in flight to finish"); got.err != nil || got.status != http.StatusOK || got.body != `[{"id":"bk-4801"}]` {
		t.Errorf("GET /bookings in flight at the stop = %d %q, %v, want 200 and the bookings", got.status, got.body, got.err)
	}
	if err := wait(t, stopped, "stop to return"); err != nil {
		t.Errorf("stop(srv, 1s) = %v, want nil once the request has finished", err)
	}
}

func TestStopRefusesNewCalls(t *testing.T) {
	srv, _, addr := serve(t)
	if err := stop(srv, time.Second); err != nil {
		t.Fatalf("stop(srv, 1s) with nothing in flight = %v, want nil", err)
	}
	if !refused(addr) { // Serve may still be starting, and closes the listener once it sees the stop
		t.Error("a new connection after stop was taken, want it refused")
	}
}

func TestStopGivesUpAfterWithin(t *testing.T) {
	srv, h, addr := serve(t)
	t.Cleanup(func() { close(h.release) })
	inFlight(addr)
	wait(t, h.entered, "the request to reach the handler")
	start := time.Now()
	err := stop(srv, 50*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("stop(srv, 50ms) with a request that outlives it = %v, want context.DeadlineExceeded", err)
	}
	if took := time.Since(start); took > time.Second {
		t.Errorf("stop(srv, 50ms) took %v, want it to give up after 50ms", took)
	}
}

func TestStopWithNoLimit(t *testing.T) {
	srv, h, addr := serve(t)
	res := inFlight(addr)
	wait(t, h.entered, "the request to reach the handler")
	stopped := make(chan error, 1)
	go func() { stopped <- stop(srv, 0) }()
	if !refused(addr) {
		t.Fatal("stop(srv, 0): new connections were still taken a second later, want them refused at once")
	}
	close(h.release)
	if got := wait(t, res, "the request in flight to finish"); got.err != nil || got.status != http.StatusOK {
		t.Errorf("GET /bookings in flight at the stop = %d, %v, want 200", got.status, got.err)
	}
	if err := wait(t, stopped, "stop to return"); err != nil {
		t.Errorf("stop(srv, 0) = %v, want nil: a zero within waits for every request in flight", err)
	}
}

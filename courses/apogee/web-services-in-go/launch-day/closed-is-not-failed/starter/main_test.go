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

// start serves a health check on a new 127.0.0.1 listener through serve, and passes on what serve returns.
func start(t *testing.T) (*http.Server, net.Listener, <-chan error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"status":"ok"}`)
	})}
	t.Cleanup(func() { srv.Close() })
	done := make(chan error, 1)
	go func() { done <- serve(srv, ln) }()
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Get("http://" + ln.Addr().String() + "/health")
	if err != nil {
		t.Fatalf("GET /health while serving: %v", err)
	}
	resp.Body.Close()
	return srv, ln, done
}

// returned waits up to a second for serve to return.
func returned(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(time.Second):
		t.Fatal("serve hadn't returned a second after the server stopped")
		return nil
	}
}

func TestServeStoppedOnPurpose(t *testing.T) {
	tests := []struct {
		name, call string
		stop       func(srv *http.Server) error
	}{
		{"the deploy tool stops it", "Shutdown", func(srv *http.Server) error { return srv.Shutdown(context.Background()) }},
		{"someone stops it in a hurry", "Close", func(srv *http.Server) error { return srv.Close() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _, done := start(t)
			if err := tt.stop(srv); err != nil {
				t.Fatal(err)
			}
			if err := returned(t, done); err != nil {
				t.Errorf("serve after %s = %v, want nil: the server was stopped on purpose", tt.call, err)
			}
		})
	}
}

func TestServeFails(t *testing.T) {
	_, ln, done := start(t)
	ln.Close() // the listener breaks underneath the server
	err := returned(t, done)
	if !errors.Is(err, net.ErrClosed) {
		t.Errorf("serve with its listener closed underneath = %v, want that failure, net.ErrClosed", err)
	}
}

func TestServeClosedListener(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ln.Close()
	srv := &http.Server{Handler: http.NotFoundHandler()}
	defer srv.Close()
	if err := serve(srv, ln); err == nil {
		t.Error("serve on a listener already closed = nil, want the failure")
	}
}

package main

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"testing/fstest"
	"time"
)

// started is what start returned.
type started struct {
	addr string
	stop func() error
	err  error
}

// within runs f and gives up waiting for it after a second.
func within(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("%s = no return after 1s, want it to return at once", what)
	}
}

func TestStart(t *testing.T) {
	s := &Service{Stations: stationTable(), Guide: fstest.MapFS{"index.html": {Data: []byte("<h1>Guide</h1>\n")}}}
	var got started
	within(t, "start(s)", func() { got.addr, got.stop, got.err = start(s) })
	if got.err != nil {
		t.Fatalf("start(s) error = %v, want nil", got.err)
	}
	if got.stop == nil {
		t.Fatalf("start(s) stop = nil, want the function stopFunc builds")
	}
	t.Cleanup(func() { // stops the server if the test ends early; stopping twice does nothing
		done := make(chan struct{})
		go func() { got.stop(); close(done) }()
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	})
	host, port, err := net.SplitHostPort(got.addr)
	if err != nil || host != "127.0.0.1" || port == "0" {
		t.Fatalf("start(s) addr = %q, want 127.0.0.1 and the port the system picked", got.addr)
	}

	client := &http.Client{Timeout: 2 * time.Second, Transport: &http.Transport{}}
	defer client.CloseIdleConnections()
	for _, path := range []string{"/stations", "/health"} {
		resp, err := client.Get("http://" + got.addr + path)
		if err != nil {
			t.Fatalf("GET %s error = %v, want an answer from the running service", path, err)
		}
		var body any
		decodeErr := json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK || decodeErr != nil {
			t.Errorf("GET %s = %s, want 200 and JSON", path, resp.Status)
		}
	}

	var stopErr error
	within(t, "stop()", func() { stopErr = got.stop() })
	if stopErr != nil {
		t.Errorf("stop() = %v, want nil", stopErr)
	}
	fresh := &http.Client{Timeout: 2 * time.Second, Transport: &http.Transport{}}
	if resp, err := fresh.Get("http://" + got.addr + "/health"); err == nil {
		resp.Body.Close()
		t.Errorf("GET /health after stop() = %s, want an error: nothing should listen there now", resp.Status)
	}
}

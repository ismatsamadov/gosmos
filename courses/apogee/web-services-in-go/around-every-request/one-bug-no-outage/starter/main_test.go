package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// syncBuffer is a buffer that the server's goroutines and the test can share.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write adds p to the buffer.
func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// String returns everything written so far.
func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func TestHandler(t *testing.T) {
	var logged syncBuffer
	s := &Service{Stations: stationTable(), Catalog: catalog(), Now: func() time.Time { return time.Date(2027, 12, 10, 15, 0, 0, 0, time.UTC) }}
	srv := httptest.NewUnstartedServer(s.handler(log.New(&logged, "", 0)))
	srv.Config.ErrorLog = log.New(io.Discard, "", 0) // the server's own report of a panic
	srv.Start()
	defer srv.Close()

	requests := []struct {
		name, path string
		status     int
	}{
		{"Ingrid's monitor", "/health", http.StatusOK},
		{"the request that broke on Monday", "/passes/cirrine-1/YST/1136", http.StatusInternalServerError},
		{"the pass Cirrine meant", "/passes/cirrine-1/YST/11364", http.StatusOK},
		{"a station that doesn't exist", "/stations/ZZZ", http.StatusNotFound},
	}
	var want []string // the log's lines, built as the answers come back
	for i, rq := range requests {
		id := fmt.Sprintf("req-%06d", i+1)
		res, err := srv.Client().Get(srv.URL + rq.path)
		if err != nil {
			t.Fatalf("GET %s (%s) error = %v, want an answer", rq.path, rq.name, err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("GET %s: reading the body: %v", rq.path, err)
		}
		if res.StatusCode != rq.status {
			t.Errorf("GET %s = %d, want %d", rq.path, res.StatusCode, rq.status)
		}
		if got := res.Header.Get("X-Request-Id"); got != id {
			t.Errorf("GET %s: X-Request-Id = %q, want %q", rq.path, got, id)
		}
		if got := res.Header.Get("Cache-Control"); got != "no-store" {
			t.Errorf("GET %s: Cache-Control = %q, want %q", rq.path, got, "no-store")
		}
		if rq.status >= 400 {
			var e ErrorBody
			if err := json.Unmarshal(body, &e); err != nil || e.RequestID != id {
				t.Errorf("GET %s: request_id in the body %q = %q, want %q", rq.path, body, e.RequestID, id)
			}
		}
		if rq.status == http.StatusInternalServerError {
			want = append(want, id+" panic: runtime error: slice bounds out of range [:5] with length 4")
		}
		want = append(want, fmt.Sprintf("%s GET %s %d %dB 0s", id, rq.path, rq.status, len(body)))
	}

	srv.Close() // waits for every handler, so the log is complete
	got := strings.Split(strings.TrimSuffix(logged.String(), "\n"), "\n")
	for i, line := range want {
		switch {
		case i >= len(got):
			t.Errorf("log line %d is missing, want %q", i+1, line)
		case got[i] != line:
			t.Errorf("log line %d = %q, want %q", i+1, got[i], line)
		}
	}
	if len(got) > len(want) {
		t.Errorf("log has %d lines, want %d: %q", len(got), len(want), got)
	}
}

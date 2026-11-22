package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// cirrine is the pass ID of Cirrine-1's pass of 8 November.
const cirrine = "cirrine-1/LAN/16402"

// lockedBuffer is a buffer the station's handler writes to while the test reads it.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write adds p to the buffer.
func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// String returns everything written so far.
func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// station is a stand-in for Lantern Atoll's uplinkd: it notes the Apogee-Pass header of each request,
// writes "lines queued" with the request's context to log, and answers 202.
type station struct {
	log     *slog.Logger
	mu      sync.Mutex
	headers []string
}

// ServeHTTP notes r's header, logs the queued lines and answers 202 Accepted.
func (s *station) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.headers = append(s.headers, r.Header.Get("Apogee-Pass"))
	s.mu.Unlock()
	s.log.InfoContext(r.Context(), "lines queued", "lines", 1)
	w.WriteHeader(http.StatusAccepted)
}

// lastHeader returns the Apogee-Pass header of the last request the station got.
func (s *station) lastHeader() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.headers) == 0 {
		return "(no request)"
	}
	return s.headers[len(s.headers)-1]
}

// newStation starts the stand-in behind handler wrap(st) and returns it with its log and server.
func newStation(t *testing.T, wrap func(http.Handler) http.Handler) (*station, *lockedBuffer, *httptest.Server) {
	t.Helper()
	var log lockedBuffer
	st := &station{log: slog.New(passHandler{slog.NewJSONHandler(&log, &slog.HandlerOptions{ReplaceAttr: dropTime})})}
	srv := httptest.NewServer(wrap(st))
	t.Cleanup(srv.Close)
	return st, &log, srv
}

// post sends POST /uplink to srv with ctx through client, and returns the request it sent.
func post(t *testing.T, client *http.Client, srv *httptest.Server, ctx context.Context) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, srv.URL+"/uplink", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST /uplink: %v", err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("POST /uplink = %d, want 202", res.StatusCode)
	}
	return req
}

// noWrap leaves the station's handler as it is.
func noWrap(h http.Handler) http.Handler { return h }

// passClient returns a client whose transport is passTransport around srv's own.
func passClient(srv *httptest.Server) *http.Client {
	return &http.Client{Transport: passTransport{base: srv.Client().Transport}, Timeout: 2 * time.Second}
}

func TestRoundTrip(t *testing.T) {
	t.Run("the pass goes out in Apogee-Pass", func(t *testing.T) {
		st, _, srv := newStation(t, noWrap)
		req := post(t, passClient(srv), srv, withPass(context.Background(), cirrine))
		if got := st.lastHeader(); got != cirrine {
			t.Errorf("POST /uplink with %s in its context: the station got Apogee-Pass %q, want %q", cirrine, got, cirrine)
		}
		if got := req.Header.Get("Apogee-Pass"); got != "" {
			t.Errorf("after POST /uplink, the caller's request has Apogee-Pass %q, want it unchanged, with none", got)
		}
	})
	t.Run("no pass, no header", func(t *testing.T) {
		st, _, srv := newStation(t, noWrap)
		post(t, passClient(srv), srv, context.Background())
		if got := st.lastHeader(); got != "" {
			t.Errorf("POST /uplink with no pass in its context: the station got Apogee-Pass %q, want none", got)
		}
	})
}

func TestAcceptPass(t *testing.T) {
	tests := []struct {
		name, header, want string
	}{
		{"a header puts the pass in the context", cirrine, cirrine},
		{"an empty header is ignored", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			h := acceptPass(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got, _ = passFrom(r.Context())
			}))
			req := httptest.NewRequest(http.MethodPost, "/uplink", nil)
			req.Header["Apogee-Pass"] = []string{tt.header}
			h.ServeHTTP(httptest.NewRecorder(), req)
			if got != tt.want {
				t.Errorf("acceptPass with Apogee-Pass %q: the handler's context holds %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}

// TestStationLine sends a pass across, as headquarters does, and reads the station's line.
func TestStationLine(t *testing.T) {
	_, log, srv := newStation(t, acceptPass)
	post(t, passClient(srv), srv, withPass(context.Background(), cirrine))
	var rec map[string]any
	raw := strings.TrimSpace(log.String())
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		t.Fatalf("the station wrote %q, want one JSON line", raw)
	}
	if rec["pass_id"] != cirrine {
		t.Errorf("headquarters sent %s, and the station wrote %s, want its pass_id", cirrine, raw)
	}
}

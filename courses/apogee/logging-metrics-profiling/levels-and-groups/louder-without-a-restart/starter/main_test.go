package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// send makes one request to /debug/level on the debug port srv and returns the status and the answer's text.
func send(t *testing.T, srv *httptest.Server, method, body string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(method, srv.URL+"/debug/level", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	client := srv.Client()
	client.Timeout = 2 * time.Second
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s /debug/level %q: %v", method, body, err)
	}
	defer res.Body.Close()
	text, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s /debug/level %q: reading the answer: %v", method, body, err)
	}
	return res.StatusCode, strings.TrimSpace(string(text))
}

func TestLevelHandler(t *testing.T) {
	tests := []struct {
		name         string
		method, body string
		from         slog.Level // the level before the request
		status       int
		text         string     // the answer's text, checked for a 400
		level        slog.Level // the level after the request
	}{
		{"the night shift turns LAN up", "PUT", "debug", slog.LevelInfo, http.StatusNoContent, "", slog.LevelDebug},
		{"and down again, in capitals", "PUT", "WARN", slog.LevelDebug, http.StatusNoContent, "", slog.LevelWarn},
		{"a level nobody has", "PUT", "loud", slog.LevelDebug, http.StatusBadRequest, `unknown level "loud"`, slog.LevelDebug},
		{"a GET", "GET", "", slog.LevelDebug, http.StatusMethodNotAllowed, "", slog.LevelDebug},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v slog.LevelVar
			v.Set(tt.from)
			srv := httptest.NewServer(levelHandler(&v))
			defer srv.Close()
			call := tt.method + " /debug/level"
			if tt.body != "" {
				call += " " + tt.body
			}
			status, text := send(t, srv, tt.method, tt.body)
			if status != tt.status {
				t.Errorf("%s = %d, want %d", call, status, tt.status)
			} else if tt.status == http.StatusBadRequest && text != tt.text {
				t.Errorf("%s answered %q, want %q", call, text, tt.text)
			}
			if got := v.Level(); got != tt.level {
				t.Errorf("after %s from %v: level = %v, want %v", call, tt.from, got, tt.level)
			}
		})
	}
}

// TestLevelHandlerLive turns a running logger up, as the night shift turns up Lantern Atoll.
func TestLevelHandlerLive(t *testing.T) {
	var v slog.LevelVar
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: &v, ReplaceAttr: dropTime}))
	srv := httptest.NewServer(levelHandler(&v))
	defer srv.Close()

	l.Debug("frame received", "bytes", 128)
	if buf.Len() != 0 {
		t.Fatalf("a DEBUG record at INFO wrote %q, want nothing", buf.String())
	}
	send(t, srv, "PUT", "debug")
	l.Debug("frame received", "bytes", 128)
	if !strings.Contains(buf.String(), `"msg":"frame received"`) {
		t.Errorf("after PUT /debug/level debug, a DEBUG record wrote %q, want its line", buf.String())
	}
}

// answers is a response that notes every status a handler writes to it.
type answers struct {
	*httptest.ResponseRecorder
	statuses []int
}

// WriteHeader notes status and passes it on.
func (a *answers) WriteHeader(status int) {
	a.statuses = append(a.statuses, status)
	a.ResponseRecorder.WriteHeader(status)
}

// TestLevelHandlerAnswersOnce checks that a level nobody has gets its 400 and nothing after it.
func TestLevelHandlerAnswersOnce(t *testing.T) {
	var v slog.LevelVar
	a := &answers{ResponseRecorder: httptest.NewRecorder()}
	levelHandler(&v).ServeHTTP(a, httptest.NewRequest(http.MethodPut, "/debug/level", strings.NewReader("loud")))
	if len(a.statuses) != 1 || a.statuses[0] != http.StatusBadRequest {
		t.Errorf("PUT /debug/level loud wrote the statuses %v, want [400] alone: nothing may follow the 400", a.statuses)
	}
}

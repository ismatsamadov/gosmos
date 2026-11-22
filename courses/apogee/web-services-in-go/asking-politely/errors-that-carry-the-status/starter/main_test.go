package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// startWeather runs a test weather service: Cape Wren's silent mast, a fault at Lantern Atoll
// explained in a whole megabyte, Ystmark stalled behind its slow link, and 404 for any code it
// doesn't know.
func startWeather(t *testing.T) *httptest.Server {
	t.Helper()
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{code}", func(w http.ResponseWriter, r *http.Request) {
		switch code := r.PathValue("code"); code {
		case "CWR":
			w.Header().Set("Retry-After", "60")
			http.Error(w, "mast silent since 2027-11-08T12:51:00Z", http.StatusServiceUnavailable)
		case "LAN":
			http.Error(w, strings.Repeat("x", 1<<20), http.StatusInternalServerError)
		case "YST":
			select {
			case <-release:
			case <-r.Context().Done():
			}
		default:
			http.Error(w, fmt.Sprintf("no station %q", code), http.StatusNotFound)
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) })
	return srv
}

func TestReadingStatusError(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 2*time.Second)
	tests := []struct {
		name       string
		code       string
		wantCode   int
		wantStatus string
		wantBody   string
		wantMsg    string
	}{
		{"Cape Wren's silent mast", "CWR", 503, "503 Service Unavailable", "mast silent since 2027-11-08T12:51:00Z",
			"weather for CWR: 503 Service Unavailable: mast silent since 2027-11-08T12:51:00Z"},
		{"a typo for Cape Wren", "CRW", 404, "404 Not Found", `no station "CRW"`,
			`weather for CRW: 404 Not Found: no station "CRW"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.Reading(context.Background(), tt.code)
			se, ok := errors.AsType[*StatusError](err)
			if !ok {
				t.Fatalf("errors.AsType[*StatusError](Reading(ctx, %q)) found none in %v, want the %d", tt.code, err, tt.wantCode)
			}
			if se.StatusCode != tt.wantCode || se.Status != tt.wantStatus || se.Body != tt.wantBody {
				t.Errorf("Reading(ctx, %q) StatusError = %+v, want {StatusCode:%d Status:%s Body:%s}",
					tt.code, *se, tt.wantCode, tt.wantStatus, tt.wantBody)
			}
			if err.Error() != tt.wantMsg {
				t.Errorf("Reading(ctx, %q) error = %q, want %q", tt.code, err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestReadingKeepsOnlyTheStart(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 2*time.Second)
	_, err := c.Reading(context.Background(), "LAN")
	se, ok := errors.AsType[*StatusError](err)
	if !ok {
		t.Fatalf(`errors.AsType[*StatusError](Reading(ctx, "LAN")) found none in %.80v, want the 500`, err)
	}
	if se.StatusCode != http.StatusInternalServerError {
		t.Errorf(`Reading(ctx, "LAN") StatusError.StatusCode = %d, want 500`, se.StatusCode)
	}
	if len(se.Body) != 512 {
		t.Errorf(`Reading(ctx, "LAN") kept %d bytes of a 1 MB body, want 512`, len(se.Body))
	}
}

func TestReadingTimeoutIsNoStatusError(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 50*time.Millisecond)
	done := make(chan error, 1)
	go func() {
		_, err := c.Reading(context.Background(), "YST")
		done <- err
	}()
	var err error
	select {
	case err = <-done:
	case <-time.After(time.Second):
		t.Fatalf(`Reading(ctx, "YST") = no answer after 1s, want a timeout error after 50ms`)
	}
	if _, ok := errors.AsType[*StatusError](err); ok {
		t.Errorf(`Reading(ctx, "YST") error = %v, a StatusError, want a timeout: no status ever came`, err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf(`Reading(ctx, "YST") error = %v, want one that errors.Is context.DeadlineExceeded`, err)
	}
}

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// startWeather runs a test weather service: a reading for Port Ellery, Cape Wren's silent
// mast, half a reading from Lantern Atoll, and 404 for any code it doesn't know.
func startWeather(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{code}", func(w http.ResponseWriter, r *http.Request) {
		switch code := r.PathValue("code"); code {
		case "PEL":
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"station":"PEL","at":"2027-11-10T14:55:00Z","wind_kmh":26,"gust_kmh":35,"temp_c":10.0}`+"\n")
		case "CWR":
			w.Header().Set("Retry-After", "60")
			http.Error(w, "mast silent since 2027-11-08T12:51:00Z", http.StatusServiceUnavailable)
		case "LAN":
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"station":"LAN","at":`)
		default:
			http.Error(w, fmt.Sprintf("no station %q", code), http.StatusNotFound)
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// bodyCounter is a transport that counts the bodies it hands out and the ones closed.
type bodyCounter struct {
	next           http.RoundTripper
	mu             sync.Mutex
	opened, closed int
}

// RoundTrip sends req and wraps the answer's body so that closing it is counted.
func (b *bodyCounter) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := b.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b.mu.Lock()
	b.opened++
	b.mu.Unlock()
	resp.Body = &countedBody{ReadCloser: resp.Body, counter: b}
	return resp, nil
}

// open returns how many bodies are still open.
func (b *bodyCounter) open() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.opened - b.closed
}

// countedBody is a body whose first Close is counted.
type countedBody struct {
	io.ReadCloser
	counter *bodyCounter
	once    sync.Once
}

// Close counts the close, once, and closes the body underneath.
func (c *countedBody) Close() error {
	c.once.Do(func() {
		c.counter.mu.Lock()
		c.counter.closed++
		c.counter.mu.Unlock()
	})
	return c.ReadCloser.Close()
}

func TestReadingClosesEveryBody(t *testing.T) {
	srv := startWeather(t)
	counter := &bodyCounter{next: srv.Client().Transport}
	c := &WeatherClient{BaseURL: srv.URL, HTTP: &http.Client{Timeout: 2 * time.Second, Transport: counter}}
	tests := []struct {
		name string
		code string
	}{
		{"Port Ellery's reading, a 200", "PEL"},
		{"a typo for Cape Wren, a 404", "CRW"},
		{"Cape Wren's silent mast, a 503", "CWR"},
		{"half a reading from Lantern Atoll, a 200 with broken JSON", "LAN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := counter.open()
			c.Reading(context.Background(), tt.code)
			if open := counter.open() - before; open != 0 {
				t.Errorf("Reading(ctx, %q) left %d body open, want 0", tt.code, open)
			}
		})
	}
}

func TestReadingStillReads(t *testing.T) {
	srv := startWeather(t)
	c := &WeatherClient{BaseURL: srv.URL, HTTP: &http.Client{Timeout: 2 * time.Second, Transport: srv.Client().Transport}}
	w, err := c.Reading(context.Background(), "PEL")
	if err != nil || w.Station != "PEL" || w.WindKmh != 26 {
		t.Errorf(`Reading(ctx, "PEL") = %+v, %v, want PEL's reading of 26 km/h and nil`, w, err)
	}
}

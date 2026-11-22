package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// startWeather runs a test weather service: Port Ellery answers at once, and Cape Wren
// stalls as it did last Tuesday, until the caller or the test gives up.
func startWeather(t *testing.T) *httptest.Server {
	t.Helper()
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/PEL", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"station":"PEL","at":"2027-11-08T08:55:00Z","wind_kmh":22,"gust_kmh":31,"temp_c":9.5}`+"\n")
	})
	mux.HandleFunc("GET /weather/CWR", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) })
	return srv
}

func TestNewWeatherClient(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 50*time.Millisecond)
	if c == nil || c.HTTP == nil {
		t.Fatalf("newWeatherClient(url, 50ms) = %+v, want a WeatherClient with an HTTP client", c)
	}
	if c.BaseURL != srv.URL {
		t.Errorf("newWeatherClient(url, 50ms).BaseURL = %q, want %q", c.BaseURL, srv.URL)
	}
	if c.HTTP == http.DefaultClient {
		t.Errorf("newWeatherClient(url, 50ms).HTTP = http.DefaultClient, want a client of its own")
	}
	if c.HTTP.Timeout != 50*time.Millisecond {
		t.Errorf("newWeatherClient(url, 50ms).HTTP.Timeout = %v, want 50ms", c.HTTP.Timeout)
	}
	if http.DefaultClient.Timeout != 0 {
		t.Errorf("http.DefaultClient.Timeout = %v, want 0: every other caller shares that client", http.DefaultClient.Timeout)
	}
}

func TestNewWeatherClientGivesUp(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 50*time.Millisecond)
	done := make(chan error, 1)
	go func() {
		_, err := c.Reading(context.Background(), "CWR")
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatalf(`Reading(ctx, "CWR") error = nil, want a timeout`)
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf(`Reading(ctx, "CWR") error = %v, want one that errors.Is context.DeadlineExceeded`, err)
		}
	case <-time.After(time.Second):
		t.Fatalf(`Reading(ctx, "CWR") = no answer after 1s, want a timeout error after 50ms`)
	}
}

func TestNewWeatherClientCalm(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 2*time.Second)
	w, err := c.Reading(context.Background(), "PEL")
	if err != nil {
		t.Fatalf(`Reading(ctx, "PEL") error = %v, want nil`, err)
	}
	at := time.Date(2027, 11, 8, 8, 55, 0, 0, time.UTC)
	if w.Station != "PEL" || w.WindKmh != 22 || !w.At.Equal(at) {
		t.Errorf(`Reading(ctx, "PEL") = %+v, want PEL's reading of 22 km/h at 08:55`, w)
	}
}

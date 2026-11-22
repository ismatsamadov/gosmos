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
// stalls until the caller or the test gives up.
func startWeather(t *testing.T) *httptest.Server {
	t.Helper()
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/PEL", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"station":"PEL","at":"2027-11-09T08:55:00Z","wind_kmh":19,"gust_kmh":27,"temp_c":8.5}`+"\n")
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

// reading calls c.Reading in a goroutine and gives up waiting for it after a second.
func reading(t *testing.T, ctx context.Context, c *WeatherClient, code string) (Weather, error) {
	t.Helper()
	type result struct {
		w   Weather
		err error
	}
	done := make(chan result, 1)
	go func() {
		w, err := c.Reading(ctx, code)
		done <- result{w, err}
	}()
	select {
	case r := <-done:
		return r.w, r.err
	case <-time.After(time.Second):
		t.Fatalf("Reading(ctx, %q) = no answer after 1s, want one as soon as ctx ends", code)
		return Weather{}, nil
	}
}

func TestReadingStopsAtTheDeadline(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 2*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := reading(t, ctx, c, "CWR")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf(`Reading(ctx ending after 50ms, "CWR") error = %v, want one that errors.Is context.DeadlineExceeded`, err)
	}
}

func TestReadingAlreadyCanceled(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 2*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	w, err := reading(t, ctx, c, "PEL")
	if !errors.Is(err, context.Canceled) {
		t.Errorf(`Reading(canceled ctx, "PEL") = %+v, %v, want an error that errors.Is context.Canceled`, w, err)
	}
}

func TestReadingCalm(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 2*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	w, err := reading(t, ctx, c, "PEL")
	if err != nil {
		t.Fatalf(`Reading(ctx, "PEL") error = %v, want nil`, err)
	}
	if w.Station != "PEL" || w.WindKmh != 19 {
		t.Errorf(`Reading(ctx, "PEL") = %+v, want PEL's reading of 19 km/h`, w)
	}
}

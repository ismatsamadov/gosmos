package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// startWeather runs a test weather service: a reading for Port Ellery, Cape Wren's silent
// mast, a fault at Lantern Atoll, and 404 for any code it doesn't know.
func startWeather(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{code}", func(w http.ResponseWriter, r *http.Request) {
		switch code := r.PathValue("code"); code {
		case "PEL":
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"station":"PEL","at":"2027-11-10T08:55:00Z","wind_kmh":24,"gust_kmh":33,"temp_c":8.0}`+"\n")
		case "CWR":
			w.Header().Set("Retry-After", "60")
			http.Error(w, "mast silent since 2027-11-08T12:51:00Z", http.StatusServiceUnavailable)
		case "LAN":
			http.Error(w, "internal error", http.StatusInternalServerError)
		default:
			http.Error(w, fmt.Sprintf("no station %q", code), http.StatusNotFound)
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestReadingChecksTheStatus(t *testing.T) {
	srv := startWeather(t)
	c := newWeatherClient(srv.URL, 2*time.Second)
	tests := []struct {
		name string
		code string
		want string // "" for a reading, else what the error must contain
	}{
		{"Port Ellery's reading", "PEL", ""},
		{"a typo for Cape Wren", "CRW", "404 Not Found"},
		{"Cape Wren's silent mast", "CWR", "503 Service Unavailable"},
		{"a fault in the weather service", "LAN", "500 Internal Server Error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := c.Reading(context.Background(), tt.code)
			if tt.want == "" {
				if err != nil || w.Station != tt.code {
					t.Errorf("Reading(ctx, %q) = %+v, %v, want %s's reading and nil", tt.code, w, err, tt.code)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Reading(ctx, %q) error = %v, want one containing %q", tt.code, err, tt.want)
			}
			if w != (Weather{}) {
				t.Errorf("Reading(ctx, %q) = %+v, want an empty Weather with the error", tt.code, w)
			}
		})
	}
}

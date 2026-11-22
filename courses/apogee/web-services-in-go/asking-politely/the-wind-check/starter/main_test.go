package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// answer is what the test weather service sends for one station.
type answer struct {
	status int     // 503 for a silent mast; otherwise 200 with the reading
	text   string  // the plain-text body of a 503
	w      Weather // the reading of a 200
	stall  bool    // answer nothing until the caller or the test gives up
}

// fakeWeather is a test weather service whose answers each case sets; a code it was never
// given gets 404, as a typo does.
type fakeWeather struct {
	mu      sync.Mutex
	answers map[string]answer
}

// set makes the service give a for code from now on.
func (f *fakeWeather) set(code string, a answer) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.answers[code] = a
}

// startWeather runs the fake weather service until the test ends.
func startWeather(t *testing.T) (*fakeWeather, *httptest.Server) {
	t.Helper()
	f := &fakeWeather{answers: map[string]answer{}}
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		f.mu.Lock()
		a, known := f.answers[code]
		f.mu.Unlock()
		switch {
		case !known:
			http.Error(w, fmt.Sprintf("no station %q", code), http.StatusNotFound)
		case a.stall:
			select {
			case <-release:
			case <-r.Context().Done():
			}
		case a.status == http.StatusServiceUnavailable:
			w.Header().Set("Retry-After", "60")
			http.Error(w, a.text, a.status)
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(a.w)
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) })
	return f, srv
}

func TestWindCall(t *testing.T) {
	now := time.Date(2027, 11, 12, 20, 0, 0, 0, time.UTC)
	fresh := now.Add(-5 * time.Minute)
	reading := func(code string, windKmh float64, at time.Time) *answer {
		return &answer{w: Weather{Station: code, At: at, WindKmh: windKmh, GustKmh: windKmh + 9, TempC: 6.5}}
	}
	tests := []struct {
		name string
		code string
		a    *answer // nil leaves the code unknown to the service
		want string
	}{
		{"Port Ellery, a light wind", "PEL", reading("PEL", 22, fresh), "GO: wind 22 km/h"},
		{"Ystmark, a windy night", "YST", reading("YST", 61, fresh), "GO: wind 61 km/h, warning"},
		{"Cape Wren, too windy to track", "CWR", reading("CWR", 74, fresh), "NO-GO: wind 74 km/h, stow"},
		{"Ystmark at 49 km/h", "YST", reading("YST", 49, fresh), "GO: wind 49 km/h"},
		{"Ystmark at exactly 50 km/h", "YST", reading("YST", 50, fresh), "GO: wind 50 km/h, warning"},
		{"Cape Wren at exactly 72 km/h", "CWR", reading("CWR", 72, fresh), "NO-GO: wind 72 km/h, stow"},
		{"Lantern Atoll, read 47 minutes ago", "LAN", reading("LAN", 18, now.Add(-47*time.Minute)), "NO-GO: reading 47m0s old"},
		{"Lantern Atoll, read exactly 15 minutes ago", "LAN", reading("LAN", 18, now.Add(-15*time.Minute)), "GO: wind 18 km/h"},
		{"Ystmark at 61 km/h, read 20 minutes ago", "YST", reading("YST", 61, now.Add(-20*time.Minute)), "NO-GO: reading 20m0s old"},
		{"Salar Quieto's silent mast", "SQT", &answer{status: 503, text: "mast silent since 2027-11-12T18:40:00Z"}, "NO-GO: mast silent"},
		{"a typo for Cape Wren", "CRW", nil, "NO-GO: no weather"},
		{"Ystmark stalled behind its slow link", "YST", &answer{stall: true}, "NO-GO: no weather"},
	}
	f, srv := startWeather(t)
	calm := newWeatherClient(srv.URL, 2*time.Second)
	quick := newWeatherClient(srv.URL, 50*time.Millisecond) // for the stalled station
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calm
			if tt.a != nil {
				f.set(tt.code, *tt.a)
				if tt.a.stall {
					c = quick
				}
			}
			done := make(chan string, 1)
			go func() { done <- windCall(context.Background(), c, tt.code, now) }()
			select {
			case got := <-done:
				if got != tt.want {
					t.Errorf("windCall(ctx, c, %q, now) = %q, want %q", tt.code, got, tt.want)
				}
			case <-time.After(time.Second):
				t.Fatalf("windCall(ctx, c, %q, now) = no answer after 1s, want %q", tt.code, tt.want)
			}
		})
	}
}

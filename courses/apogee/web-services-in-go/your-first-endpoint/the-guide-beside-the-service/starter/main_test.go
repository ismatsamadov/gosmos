package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

// guide is a small customer guide, as Yusra's pages look in a test.
func guide() fstest.MapFS {
	return fstest.MapFS{
		"index.html":    {Data: []byte("<h1>The booking service: a guide for customers</h1>\n")},
		"stations.html": {Data: []byte("<h1>Stations you can book</h1>\n")},
	}
}

// get sends a GET for target through h and returns what h answered.
func get(h http.Handler, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestRoutesGuide(t *testing.T) {
	mux := (&Service{Stations: stationTable(), Guide: guide()}).routes()
	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantBody   string // part of the page, for a 200
	}{
		{"the guide's front page", "/guide/", http.StatusOK, "a guide for customers"},
		{"the stations page", "/guide/stations.html", http.StatusOK, "Stations you can book"},
		{"a page the guide doesn't have", "/guide/missing.html", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := get(mux, tt.target)
			if rec.Code != tt.wantStatus {
				t.Fatalf("GET %s = %d, want %d", tt.target, rec.Code, tt.wantStatus)
			}
			if tt.wantStatus != http.StatusOK {
				return
			}
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("GET %s body = %q, want the page with %q", tt.target, rec.Body.String(), tt.wantBody)
			}
			if ct := rec.Result().Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
				t.Errorf("GET %s Content-Type = %q, want text/html", tt.target, ct)
			}
		})
	}
}

func TestRoutesGuideRedirect(t *testing.T) {
	mux := (&Service{Stations: stationTable(), Guide: guide()}).routes()
	rec := get(mux, "/guide")
	if rec.Code != http.StatusTemporaryRedirect || rec.Header().Get("Location") != "/guide/" {
		t.Errorf("GET /guide = %d to %q, want 307 to %q", rec.Code, rec.Header().Get("Location"), "/guide/")
	}
}

func TestRoutesStillAnswer(t *testing.T) {
	mux := (&Service{Stations: stationTable(), Guide: guide()}).routes()
	for _, target := range []string{"/stations", "/health"} {
		if rec := get(mux, target); rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", target, rec.Code)
		}
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// get sends a GET for target through h and returns what h answered.
func get(h http.Handler, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestRoutesHealth(t *testing.T) {
	mux := (&Service{Stations: stationTable()}).routes()
	rec := get(mux, "/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health = %d, want 200", rec.Code)
	}
	if ct := rec.Result().Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("GET /health Content-Type = %q, want %q", ct, "application/json")
	}
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || len(got) != 1 || got["status"] != "ok" {
		t.Errorf(`GET /health body = %q, want {"status":"ok"}`, rec.Body.String())
	}
}

func TestRoutesStations(t *testing.T) {
	mux := (&Service{Stations: stationTable()}).routes()
	rec := get(mux, "/stations")
	var got []Station
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &got) != nil || len(got) != 5 {
		t.Errorf("GET /stations = %d %.60q, want 200 and the five stations", rec.Code, rec.Body.String())
	}
}

func TestRoutesUnknownPaths(t *testing.T) {
	mux := (&Service{Stations: stationTable()}).routes()
	for _, target := range []string{"/stations/", "/statoins", "/healthz", "/"} {
		if rec := get(mux, target); rec.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d %.40q, want 404 from the mux", target, rec.Code, rec.Body.String())
		}
	}
}

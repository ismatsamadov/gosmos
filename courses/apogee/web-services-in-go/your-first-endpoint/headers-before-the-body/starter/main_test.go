package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSONCreated(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusCreated, stationTable()[1])
	res := rec.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("writeJSON(rec, 201, PEL) status = %d, want 201", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("writeJSON(rec, 201, PEL) Content-Type = %q, want %q", ct, "application/json")
	}
	body := rec.Body.String()
	var got Station
	if err := json.Unmarshal([]byte(body), &got); err != nil || got.Code != "PEL" {
		t.Errorf("writeJSON(rec, 201, PEL) body = %q, want Port Ellery as JSON", body)
	}
	if !strings.HasSuffix(body, "}\n") {
		t.Errorf("writeJSON(rec, 201, PEL) body = %q, want the JSON followed by one newline", body)
	}
}

func TestWriteJSONCannotMarshal(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusOK, make(chan int))
	res := rec.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("writeJSON(rec, 200, a channel) status = %d, want 500", res.StatusCode)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "internal error" {
		t.Errorf("writeJSON(rec, 200, a channel) body = %q, want %q", got, "internal error")
	}
	if ct := res.Header.Get("Content-Type"); strings.HasPrefix(ct, "application/json") {
		t.Errorf("writeJSON(rec, 200, a channel) Content-Type = %q, want plain text, since no JSON was sent", ct)
	}
}

func TestListStationsThroughWriteJSON(t *testing.T) {
	s := &Service{Stations: stationTable()}
	rec := httptest.NewRecorder()
	s.listStations(rec, httptest.NewRequest(http.MethodGet, "/stations", nil))
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("GET /stations = %d, want 200", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("GET /stations Content-Type = %q, want %q", ct, "application/json")
	}
	var got []Station
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || len(got) != 5 {
		t.Errorf("GET /stations body = %.60q, want the five stations as JSON", rec.Body.String())
	}
}

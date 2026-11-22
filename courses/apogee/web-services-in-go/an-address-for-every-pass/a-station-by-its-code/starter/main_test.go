package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

// newTestService returns the booking service with the station table and a one-page guide.
func newTestService() *Service {
	return &Service{
		Stations: stationTable(),
		Guide:    fstest.MapFS{"index.html": {Data: []byte("<h1>Customer guide</h1>\n")}},
	}
}

// get sends a GET request for path through the service's routes.
func get(path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	newTestService().routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestGetStation(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		status int
		want   string // the station's name, or the whole plain-text answer of a 404
	}{
		{"Salar Quieto by its code", "/stations/SQT", 200, "Salar Quieto"},
		{"Cape Wren by its code", "/stations/CWR", 200, "Cape Wren"},
		{"a code no station has", "/stations/ZZZ", 404, "no station \"ZZZ\"\n"},
		{"a code in lowercase", "/stations/sqt", 404, "no station \"sqt\"\n"},
		{"a segment too many", "/stations/SQT/extra", 404, "404 page not found\n"},
		{"no code at all", "/stations/", 404, "404 page not found\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := get(tt.path)
			if tt.status != http.StatusOK {
				if rec.Code != tt.status || rec.Body.String() != tt.want {
					t.Errorf("GET %s = %d %q, want %d %q", tt.path, rec.Code, rec.Body.String(), tt.status, tt.want)
				}
				return
			}
			var st Station
			if err := json.Unmarshal(rec.Body.Bytes(), &st); rec.Code != http.StatusOK || err != nil {
				t.Fatalf("GET %s = %d %q, want 200 and a station", tt.path, rec.Code, rec.Body.String())
			}
			if st.Name != tt.want {
				t.Errorf("GET %s: name = %q, want %q", tt.path, st.Name, tt.want)
			}
		})
	}
}

func TestListStationsStillThere(t *testing.T) {
	rec := get("/stations")
	var list []Station
	if err := json.Unmarshal(rec.Body.Bytes(), &list); rec.Code != http.StatusOK || err != nil {
		t.Fatalf("GET /stations = %d %q, want 200 and a list", rec.Code, rec.Body.String())
	}
	if len(list) != 5 {
		t.Errorf("GET /stations: %d stations, want 5", len(list))
	}
}

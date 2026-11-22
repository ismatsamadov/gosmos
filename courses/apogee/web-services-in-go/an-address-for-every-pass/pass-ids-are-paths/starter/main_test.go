package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

// newTestService returns the booking service with the station table, the catalog and a one-page guide.
func newTestService() *Service {
	return &Service{
		Stations: stationTable(),
		Catalog:  catalog(),
		Guide:    fstest.MapFS{"index.html": {Data: []byte("<h1>Customer guide</h1>\n")}},
	}
}

// get sends a GET request for path through the service's routes.
func get(path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	newTestService().routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestGetPass(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		status int
		want   string // the pass ID, or the whole plain-text answer of a 404
	}{
		{"Furrow-A over Salar Quieto on 6 December", "/passes/furrow-a/SQT/15864", 200, "furrow-a/SQT/15864"},
		{"Pipit-1 over Port Ellery, orbit 08154", "/passes/pipit-1/PEL/08154", 200, "pipit-1/PEL/08154"},
		{"Furrow-B over Ystmark", "/passes/furrow-b/YST/15890", 200, "furrow-b/YST/15890"},
		{"an orbit with no pass", "/passes/furrow-a/SQT/15865", 404, "no pass furrow-a/SQT/15865\n"},
		{"an orbit without its leading zero", "/passes/pipit-1/PEL/8154", 404, "no pass pipit-1/PEL/8154\n"},
		{"a pass ID cut short", "/passes/furrow-a/SQT", 404, "404 page not found\n"},
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
			var p PassView
			if err := json.Unmarshal(rec.Body.Bytes(), &p); rec.Code != http.StatusOK || err != nil {
				t.Fatalf("GET %s = %d %q, want 200 and a pass", tt.path, rec.Code, rec.Body.String())
			}
			if p.ID != tt.want {
				t.Errorf("GET %s: id = %q, want %q", tt.path, p.ID, tt.want)
			}
		})
	}
}

func TestGetPassFields(t *testing.T) {
	const path = "/passes/furrow-a/SQT/15864"
	rec := get(path)
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); rec.Code != http.StatusOK || err != nil {
		t.Fatalf("GET %s = %d %q, want 200 and a pass", path, rec.Code, rec.Body.String())
	}
	want := map[string]any{
		"id":           "furrow-a/SQT/15864",
		"sat_id":       "furrow-a",
		"station":      "SQT",
		"aos":          "2027-12-06T14:02:10Z",
		"los":          "2027-12-06T14:11:52Z",
		"max_elev_deg": 47.9,
	}
	for key, w := range want {
		if fmt.Sprint(got[key]) != fmt.Sprint(w) {
			t.Errorf("GET %s: %q = %v, want %v", path, key, got[key], w)
		}
	}
}

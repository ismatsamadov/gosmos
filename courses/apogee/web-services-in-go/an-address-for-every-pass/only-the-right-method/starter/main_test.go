package main

import (
	"io"
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

func TestRoutesMethods(t *testing.T) {
	tests := []struct {
		name         string
		method, path string
		status       int
		allow        string
	}{
		{"the station list", "GET", "/stations", 200, ""},
		{"the health check", "GET", "/health", 200, ""},
		{"the customer guide", "GET", "/guide/", 200, ""},
		{"Furrowlight's POST by mistake", "POST", "/stations", 405, "GET, HEAD"},
		{"a DELETE of the health check", "DELETE", "/health", 405, "GET, HEAD"},
		{"a POST to the guide", "POST", "/guide/", 405, "GET, HEAD"},
		{"a path nobody registered", "GET", "/passes", 404, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestService().routes().ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))
			if rec.Code != tt.status {
				t.Errorf("%s %s = %d, want %d", tt.method, tt.path, rec.Code, tt.status)
			}
			if got := rec.Header().Get("Allow"); got != tt.allow {
				t.Errorf("%s %s: Allow = %q, want %q", tt.method, tt.path, got, tt.allow)
			}
		})
	}
}

func TestRoutesHead(t *testing.T) {
	srv := httptest.NewServer(newTestService().routes())
	defer srv.Close()
	req, err := http.NewRequest(http.MethodHead, srv.URL+"/stations", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("HEAD /stations: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("HEAD /stations: reading the body: %v", err)
	}
	if res.StatusCode != http.StatusOK || len(body) != 0 {
		t.Errorf("HEAD /stations = %d with a %d-byte body, want 200 with none", res.StatusCode, len(body))
	}
	if got := res.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("HEAD /stations: Content-Type = %q, want %q", got, "application/json")
	}
}

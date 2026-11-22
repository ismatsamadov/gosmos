package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
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

// catalogIDs returns the ID of every pass in the catalog, in its order.
func catalogIDs() []string {
	var ids []string
	for _, p := range catalog() {
		ids = append(ids, p.ID)
	}
	return ids
}

func TestListPasses(t *testing.T) {
	tests := []struct {
		name   string
		target string
		want   []string // the IDs listed, in order
	}{
		{"Furrow-A over Salar Quieto", "/passes?sat=furrow-a&station=SQT",
			[]string{"furrow-a/SQT/15697", "furrow-a/SQT/15758", "furrow-a/SQT/15818", "furrow-a/SQT/15864", "furrow-a/SQT/15910"}},
		{"every pass over Salar Quieto", "/passes?station=SQT",
			[]string{"furrow-a/SQT/15697", "furrow-a/SQT/15758", "furrow-b/SQT/15786", "furrow-a/SQT/15818", "furrow-a/SQT/15864", "furrow-b/SQT/15862", "furrow-a/SQT/15910", "furrow-b/SQT/15923"}},
		{"every pass of Brine-3", "/passes?sat=brine-3",
			[]string{"brine-3/LAN/09056", "brine-3/LAN/09168", "brine-3/LAN/09213"}},
		{"an empty value is no filter", "/passes?sat=&station=CWR",
			[]string{"furrow-b/CWR/15709", "furrow-a/CWR/15879"}},
		{"no filter at all", "/passes", catalogIDs()},
		{"Brine-3 never rises over Ystmark", "/passes?sat=brine-3&station=YST", nil},
		{"Pipit-2 has no pass in the catalog", "/passes?sat=pipit-2", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestService().routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.target, nil))
			if rec.Code != http.StatusOK {
				t.Fatalf("GET %s = %d %q, want 200", tt.target, rec.Code, rec.Body.String())
			}
			if len(tt.want) == 0 {
				if body := strings.TrimSpace(rec.Body.String()); body != "[]" {
					t.Errorf("GET %s = %q, want %q", tt.target, body, "[]")
				}
				return
			}
			var passes []PassView
			if err := json.Unmarshal(rec.Body.Bytes(), &passes); err != nil {
				t.Fatalf("GET %s = %q, want a list of passes: %v", tt.target, rec.Body.String(), err)
			}
			var got []string
			for _, p := range passes {
				got = append(got, p.ID)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("GET %s = %v, want %v", tt.target, got, tt.want)
			}
		})
	}
}

func TestListPassesRefuses(t *testing.T) {
	tests := []struct {
		name   string
		target string
		want   string
	}{
		{"Thistle re-entered in July", "/passes?sat=thistle", "unknown satellite \"thistle\""},
		{"a satellite in capitals", "/passes?sat=FURROW-A&station=SQT", "unknown satellite \"FURROW-A\""},
		{"a station nobody has", "/passes?station=ZZZ", "unknown station \"ZZZ\""},
		{"a station code in lowercase", "/passes?sat=furrow-a&station=sqt", "unknown station \"sqt\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestService().routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.target, nil))
			if got := strings.TrimSpace(rec.Body.String()); rec.Code != http.StatusBadRequest || got != tt.want {
				t.Errorf("GET %s = %d %q, want 400 %q", tt.target, rec.Code, got, tt.want)
			}
		})
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestListStationsByBand(t *testing.T) {
	all := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantCodes  []string // for a 200
		wantText   string   // for a 400
	}{
		{"X-band", "/stations?band=X", 200, []string{"YST", "SQT", "CWR"}, ""},
		{"UHF, listen-only antennas left out", "/stations?band=UHF", 200, []string{"PEL", "LAN"}, ""},
		{"S-band, everywhere", "/stations?band=S", 200, all, ""},
		{"no band", "/stations", 200, all, ""},
		{"an empty band", "/stations?band=", 200, all, ""},
		{"Furrowlight's lowercase x", "/stations?band=x", 400, nil, `unknown band "x"`},
		{"La Grande's test radio, which nobody books", "/stations?band=Ka", 400, nil, `unknown band "Ka"`},
	}
	s := &Service{Stations: stationTable()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			s.listStations(rec, httptest.NewRequest(http.MethodGet, tt.target, nil))
			if rec.Code != tt.wantStatus {
				t.Fatalf("GET %s = %d, want %d", tt.target, rec.Code, tt.wantStatus)
			}
			if tt.wantStatus != http.StatusOK {
				if got := strings.TrimSpace(rec.Body.String()); got != tt.wantText {
					t.Errorf("GET %s body = %q, want %q", tt.target, got, tt.wantText)
				}
				return
			}
			var got []Station
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("GET %s body = %.60q, want a JSON list: %v", tt.target, rec.Body.String(), err)
			}
			var codes []string
			for _, st := range got {
				codes = append(codes, st.Code)
			}
			if !slices.Equal(codes, tt.wantCodes) {
				t.Errorf("GET %s = %v, want %v", tt.target, codes, tt.wantCodes)
			}
		})
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestStationsHandler(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
		want []string
	}{
		{"the whole network", &Service{Stations: stationTable()}, []string{"YST", "PEL", "SQT", "LAN", "CWR"}},
		{"a service that knows Port Ellery alone", &Service{Stations: stationTable()[1:2]}, []string{"PEL"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tt.s.stationsHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/stations", nil))
			if rec.Code != http.StatusOK {
				t.Fatalf("GET /stations through stationsHandler() = %d, want 200", rec.Code)
			}
			var got []Station
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("GET /stations body = %.60q, want a JSON list: %v", rec.Body.String(), err)
			}
			var codes []string
			for _, st := range got {
				codes = append(codes, st.Code)
			}
			if !slices.Equal(codes, tt.want) {
				t.Errorf("GET /stations through stationsHandler() = %v, want this service's own %v", codes, tt.want)
			}
		})
	}
}

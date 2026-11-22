package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

// The list is a handler: this line doesn't build unless stationList has ServeHTTP.
var _ http.Handler = stationList{}

func TestStationListServeHTTP(t *testing.T) {
	rec := httptest.NewRecorder()
	stationList(stationTable()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/stations", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("GET /stations = %d, want 200", rec.Code)
	}
	var got []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("GET /stations body = %.60q, want a JSON list of stations: %v", rec.Body.String(), err)
	}
	var codes []string
	for _, st := range got {
		codes = append(codes, fmt.Sprint(st["code"]))
	}
	if want := []string{"YST", "PEL", "SQT", "LAN", "CWR"}; !slices.Equal(codes, want) {
		t.Fatalf("GET /stations codes = %v, want %v, in table order", codes, want)
	}
	sqt := got[2]
	if sqt["name"] != "Salar Quieto" || sqt["lat_deg"] != -23.4 || fmt.Sprint(sqt["bands"]) != "[S X]" {
		t.Errorf("GET /stations gave SQT as %v, want its name, lat_deg -23.4 and bands [S X]", sqt)
	}
}

func TestEmptyStationList(t *testing.T) {
	rec := httptest.NewRecorder()
	stationList{}.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/stations", nil))
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Errorf("GET /stations with no stations = %q, want []", got)
	}
}

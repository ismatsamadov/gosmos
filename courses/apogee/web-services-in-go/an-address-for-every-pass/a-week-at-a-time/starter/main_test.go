package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

// thursday is the service's clock in these tests: Thursday 25 November 2027, 10:00 UTC.
var thursday = time.Date(2027, 11, 25, 10, 0, 0, 0, time.UTC)

// newTestService returns the booking service with the lesson's data, its clock stopped on Thursday.
func newTestService() *Service {
	return &Service{
		Stations: stationTable(),
		Catalog:  catalog(),
		Guide:    fstest.MapFS{"index.html": {Data: []byte("<h1>Customer guide</h1>\n")}},
		Now:      func() time.Time { return thursday },
	}
}

func TestWindow(t *testing.T) {
	monday := time.Date(2027, 12, 6, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		query    url.Values
		wantFrom time.Time
		wantDays int
	}{
		{"no from or days", url.Values{}, thursday, 7},
		{"one day", url.Values{"days": {"1"}}, thursday, 1},
		{"two weeks", url.Values{"days": {"14"}}, thursday, 14},
		{"from Monday 6 December", url.Values{"from": {"2027-12-06T00:00:00Z"}}, monday, 7},
		{"Hollis's time, three days", url.Values{"from": {"2027-12-06T15:02:10+01:00"}, "days": {"3"}}, monday.Add(14*time.Hour + 2*time.Minute + 10*time.Second), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, days, err := window(tt.query, thursday)
			if err != nil || !from.Equal(tt.wantFrom) || days != tt.wantDays {
				t.Errorf("window(%q) = %v, %d, %v, want %v, %d, nil", tt.query.Encode(), from, days, err, tt.wantFrom, tt.wantDays)
			}
		})
	}
}

func TestWindowRefuses(t *testing.T) {
	const fromMsg = "from must be an RFC 3339 time"
	const daysMsg = "days must be a whole number from 1 to 14"
	tests := []struct {
		name  string
		query url.Values
		want  string
	}{
		{"no days at all", url.Values{"days": {"0"}}, daysMsg},
		{"one day too many", url.Values{"days": {"15"}}, daysMsg},
		{"a week in words", url.Values{"days": {"seven"}}, daysMsg},
		{"a week and a half", url.Values{"days": {"7.5"}}, daysMsg},
		{"days counted backward", url.Values{"days": {"-3"}}, daysMsg},
		{"a weekday", url.Values{"from": {"monday"}}, fromMsg},
		{"a date with no time", url.Values{"from": {"2027-12-06"}}, fromMsg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := window(tt.query, thursday)
			if err == nil || err.Error() != tt.want {
				t.Errorf("window(%q) error = %v, want %q", tt.query.Encode(), err, tt.want)
			}
		})
	}
}

func TestListPassesWindow(t *testing.T) {
	tests := []struct {
		name  string
		query url.Values
		want  []string // the IDs listed, in order
	}{
		{"a week from Thursday 10:00", url.Values{},
			[]string{"furrow-a/SQT/15697", "brine-3/LAN/09056", "furrow-b/CWR/15709", "cirrine-1/YST/11260", "furrow-a/SQT/15758", "pipit-1/PEL/08078", "furrow-b/SQT/15786"}},
		{"one day from Thursday 10:00", url.Values{"days": {"1"}},
			[]string{"furrow-a/SQT/15697", "brine-3/LAN/09056"}},
		{"Monday 6 December", url.Values{"from": {"2027-12-06T00:00:00Z"}, "days": {"1"}},
			[]string{"pipit-1/PEL/08154", "cirrine-1/YST/11364", "brine-3/LAN/09213", "furrow-a/SQT/15864", "furrow-b/SQT/15862"}},
		{"Furrow-A from Hollis's 15:02:10, his time", url.Values{"sat": {"furrow-a"}, "from": {"2027-12-06T15:02:10+01:00"}, "days": {"1"}},
			[]string{"furrow-a/SQT/15864", "furrow-a/CWR/15879"}},
		{"a pass at exactly from is in", url.Values{"sat": {"furrow-a"}, "station": {"SQT"}, "from": {"2027-12-06T14:02:10Z"}, "days": {"1"}},
			[]string{"furrow-a/SQT/15864"}},
		{"a pass at exactly the end is out", url.Values{"sat": {"furrow-a"}, "station": {"SQT"}, "from": {"2027-12-05T14:02:10Z"}, "days": {"1"}},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/passes?" + tt.query.Encode()
			rec := httptest.NewRecorder()
			newTestService().routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
			var passes []PassView
			if err := json.Unmarshal(rec.Body.Bytes(), &passes); rec.Code != http.StatusOK || err != nil {
				t.Fatalf("GET %s = %d %q, want 200 and a list of passes", target, rec.Code, rec.Body.String())
			}
			var got []string
			for _, p := range passes {
				got = append(got, p.ID)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("GET %s = %v, want %v", target, got, tt.want)
			}
		})
	}
	target := "/passes?days=0"
	rec := httptest.NewRecorder()
	newTestService().routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	if got := strings.TrimSpace(rec.Body.String()); rec.Code != http.StatusBadRequest || got != "days must be a whole number from 1 to 14" {
		t.Errorf("GET %s = %d %q, want 400 %q", target, rec.Code, got, "days must be a whole number from 1 to 14")
	}
}

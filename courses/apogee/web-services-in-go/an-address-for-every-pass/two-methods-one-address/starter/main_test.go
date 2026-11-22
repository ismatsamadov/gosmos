package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"
)

// newTestService returns the booking service with the lesson's data and Yusra's imported bookings, its clock stopped at now.
func newTestService(t *testing.T, now time.Time) *Service {
	t.Helper()
	store := NewMemStore()
	if err := store.load(importedBookings()); err != nil {
		t.Fatalf("loading Yusra's bookings: %v", err)
	}
	return &Service{
		Stations: stationTable(),
		Catalog:  catalog(),
		Guide:    fstest.MapFS{"index.html": {Data: []byte("<h1>Customer guide</h1>\n")}},
		Store:    store,
		Now:      func() time.Time { return now },
	}
}

func TestCancelBooking(t *testing.T) {
	s := newTestService(t, time.Date(2027, 11, 25, 10, 0, 0, 0, time.UTC))
	srv := httptest.NewServer(s.routes())
	defer srv.Close()
	// The steps run in order: each one sees what the ones before it did.
	steps := []struct {
		name         string
		method, path string
		status       int
		body         string // the whole body, or "" for none; a 200 is checked as JSON instead
		allow        string
	}{
		{"Brinewatch cancels bk-4496", "DELETE", "/bookings/bk-4496", 204, "", ""},
		{"bk-4496 is gone", "GET", "/bookings/bk-4496", 404, "booking bk-4496: not found\n", ""},
		{"a second cancellation", "DELETE", "/bookings/bk-4496", 404, "booking bk-4496: not found\n", ""},
		{"bk-4497 is still there", "GET", "/bookings/bk-4497", 200, "", ""},
		{"a PUT on a booking", "PUT", "/bookings/bk-4497", 405, "Method Not Allowed\n", "DELETE, GET, HEAD"},
	}
	for _, st := range steps {
		t.Run(st.name, func(t *testing.T) {
			req, err := http.NewRequest(st.method, srv.URL+st.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			res, err := srv.Client().Do(req)
			if err != nil {
				t.Fatalf("%s %s: %v", st.method, st.path, err)
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("%s %s: reading the body: %v", st.method, st.path, err)
			}
			if res.StatusCode != st.status {
				t.Errorf("%s %s = %d %q, want %d", st.method, st.path, res.StatusCode, body, st.status)
			}
			if st.status == http.StatusOK {
				var b Booking
				if err := json.Unmarshal(body, &b); err != nil || b.ID != "bk-4497" {
					t.Errorf("%s %s = %q, want booking bk-4497", st.method, st.path, body)
				}
			} else if string(body) != st.body {
				t.Errorf("%s %s: body = %q, want %q", st.method, st.path, body, st.body)
			}
			if got := res.Header.Get("Allow"); got != st.allow {
				t.Errorf("%s %s: Allow = %q, want %q", st.method, st.path, got, st.allow)
			}
		})
	}
}

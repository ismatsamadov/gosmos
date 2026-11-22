package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

// send serves one request through h, with body and its Content-Type when one is given.
func send(h http.Handler, method, target, contentType, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// sameBooking reports whether a and b hold the same booking, times compared as instants.
func sameBooking(a, b Booking) bool {
	return a.ID == b.ID && a.SatID == b.SatID && a.Station == b.Station &&
		a.AOS.Equal(b.AOS) && a.LOS.Equal(b.LOS) && a.Note == b.Note
}

// monday is the service's clock in these tests: Monday 29 November 2027, 10:00 UTC.
var monday = time.Date(2027, 11, 29, 10, 0, 0, 0, time.UTC)

func TestCreateBooking(t *testing.T) {
	h := newTestService(t, monday).routes()
	// Apogee tries every change on a Pipit first. The bookings run in order.
	bookings := []struct {
		name, body string
		want       Booking
	}{
		{"a Pipit-1 test booking",
			`{"sat_id":"pipit-1","station":"PEL","aos":"2027-12-06T04:12:05Z","los":"2027-12-06T04:21:27Z","note":"test of POST /bookings"}`,
			Booking{ID: "bk-4501", SatID: "pipit-1", Station: "PEL", AOS: mustTime("2027-12-06T04:12:05Z"), LOS: mustTime("2027-12-06T04:21:27Z"), Note: "test of POST /bookings"}},
		{"the next test booking",
			`{"sat_id":"pipit-1","station":"PEL","aos":"2027-12-01T03:53:31Z","los":"2027-12-01T04:02:32Z"}`,
			Booking{ID: "bk-4502", SatID: "pipit-1", Station: "PEL", AOS: mustTime("2027-12-01T03:53:31Z"), LOS: mustTime("2027-12-01T04:02:32Z")}},
	}
	for _, bk := range bookings {
		t.Run(bk.name, func(t *testing.T) {
			rec := send(h, http.MethodPost, "/bookings", "application/json", bk.body)
			var got Booking
			if err := json.Unmarshal(rec.Body.Bytes(), &got); rec.Code != http.StatusCreated || err != nil {
				t.Fatalf("POST /bookings = %d %q, want 201 and the booking", rec.Code, rec.Body.String())
			}
			if !sameBooking(got, bk.want) {
				t.Errorf("POST /bookings = %+v, want %+v", got, bk.want)
			}
			if rec := send(h, http.MethodGet, "/bookings/"+bk.want.ID, "", ""); rec.Code != http.StatusOK {
				t.Errorf("GET /bookings/%s after the POST = %d %q, want 200", bk.want.ID, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestCreateBookingRefusesBrokenJSON(t *testing.T) {
	s := newTestService(t, monday)
	h := s.routes()
	bodies := []struct{ name, body string }{
		{"a body cut short", `{"sat_id":"pipit-1","station":`},
		{"an empty body", ""},
		{"a form instead of JSON", "sat_id=pipit-1&station=PEL"},
	}
	for _, bd := range bodies {
		t.Run(bd.name, func(t *testing.T) {
			rec := send(h, http.MethodPost, "/bookings", "application/json", bd.body)
			if rec.Code != http.StatusBadRequest || !strings.HasPrefix(rec.Body.String(), "request body: ") {
				t.Errorf("POST /bookings with %s = %d %q, want 400 \"request body: ...\"", bd.name, rec.Code, rec.Body.String())
			}
		})
	}
	kept, err := s.Store.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(kept) != len(importedBookings()) {
		t.Errorf("after the broken bodies the store holds %d bookings, want Yusra's %d", len(kept), len(importedBookings()))
	}
}

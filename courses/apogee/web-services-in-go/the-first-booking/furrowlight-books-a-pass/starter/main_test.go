package main

import (
	"context"
	"encoding/json"
	"fmt"
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

// friday is the service's clock in these tests: Friday 3 December 2027, 14:30 UTC.
var friday = time.Date(2027, 12, 3, 14, 30, 0, 0, time.UTC)

// hollis is the body of Furrowlight's first booking: Furrow-A over Salar Quieto on Monday 6 December.
const hollis = `{"sat_id":"furrow-a","station":"SQT","aos":"2027-12-06T14:02:10Z","los":"2027-12-06T14:11:52Z","note":"PRIORITY"}`

// bookingBody returns the JSON body that books the catalog pass p.
func bookingBody(p PassView) string {
	return `{"sat_id":"` + p.SatID + `","station":"` + p.Station + `","aos":"` + p.AOS.Format(time.RFC3339) + `","los":"` + p.LOS.Format(time.RFC3339) + `"}`
}

// errorOf decodes the ErrorBody of an answer, or returns an empty one when it isn't one.
func errorOf(rec *httptest.ResponseRecorder) ErrorBody {
	var e ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil {
		return ErrorBody{}
	}
	return e
}

func TestCreateBookingFirst(t *testing.T) {
	h := newTestService(t, friday).routes()
	rec := send(h, http.MethodPost, "/bookings", "application/json", hollis)
	var got Booking
	if err := json.Unmarshal(rec.Body.Bytes(), &got); rec.Code != http.StatusCreated || err != nil {
		t.Fatalf("POST /bookings with Hollis's booking = %d %q, want 201 and the booking", rec.Code, rec.Body.String())
	}
	if got.ID != "bk-4501" || got.SatID != "furrow-a" || got.Station != "SQT" || !got.AOS.Equal(mustTime("2027-12-06T14:02:10Z")) || got.Note != "PRIORITY" {
		t.Errorf("POST /bookings with Hollis's booking = %+v, want bk-4501, Furrow-A over SQT at 14:02:10, PRIORITY", got)
	}
	if loc := rec.Header().Get("Location"); loc != "/bookings/bk-4501" {
		t.Errorf("POST /bookings with Hollis's booking: Location = %q, want %q", loc, "/bookings/bk-4501")
	}
	if rec := send(h, http.MethodGet, "/bookings/bk-4501", "", ""); rec.Code != http.StatusOK {
		t.Errorf("GET /bookings/bk-4501 after the POST = %d %q, want 200", rec.Code, rec.Body.String())
	}
	rec = send(h, http.MethodPost, "/bookings", "application/json", hollis)
	want := ErrorBody{Code: "conflict", Message: "furrow-a over SQT at 2027-12-06T14:02:10Z: already booked"}
	if got := errorOf(rec); rec.Code != http.StatusConflict || got != want {
		t.Errorf("POST /bookings with the same pass again = %d %q, want 409 %+v", rec.Code, rec.Body.String(), want)
	}
}

func TestCreateBookingRefusals(t *testing.T) {
	s := newTestService(t, friday)
	h := s.routes()
	tests := []struct {
		name        string
		contentType string
		body        string
		status      int
		want        ErrorBody
	}{
		{"a form, as Brinewatch sent", "application/x-www-form-urlencoded", "sat_id=furrow-a&station=SQT",
			415, ErrorBody{Code: "unsupported_media_type", Message: "send the body as application/json"}},
		{"a 70,000-letter note", "application/json", strings.Replace(hollis, "PRIORITY", strings.Repeat("A", 70000), 1),
			413, ErrorBody{Code: "too_large", Message: "request body over 65536 bytes"}},
		{"sat instead of sat_id", "application/json", strings.Replace(hollis, `"sat_id"`, `"sat"`, 1),
			400, ErrorBody{Code: "bad_request", Message: `request body: json: unknown field "sat"`}},
		{"no los", "application/json", strings.Replace(hollis, `"los":"2027-12-06T14:11:52Z",`, "", 1),
			422, ErrorBody{Code: "invalid_booking", Message: "aos and los are required"}},
		{"an ID chosen by the caller", "application/json", strings.Replace(hollis, "{", `{"id":"bk-4501",`, 1),
			422, ErrorBody{Code: "invalid_booking", Message: "id is chosen by the service"}},
		{"Thistle", "application/json", strings.Replace(hollis, "furrow-a", "thistle", 1),
			422, ErrorBody{Code: "invalid_booking", Message: `unknown satellite "thistle"`}},
		{"the pass that began an hour ago", "application/json", `{"sat_id":"furrow-a","station":"SQT","aos":"2027-12-03T13:28:42Z","los":"2027-12-03T13:36:51Z"}`,
			422, ErrorBody{Code: "invalid_booking", Message: "bookings for this pass closed at 2027-12-03T13:18:42Z"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := send(h, http.MethodPost, "/bookings", tt.contentType, tt.body)
			if got := errorOf(rec); rec.Code != tt.status || got != tt.want {
				t.Errorf("POST /bookings with %s = %d %q, want %d %+v", tt.name, rec.Code, rec.Body.String(), tt.status, tt.want)
			}
		})
	}
	kept, err := s.Store.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(kept) != len(importedBookings()) {
		t.Errorf("after the refusals the store holds %d bookings, want Yusra's %d", len(kept), len(importedBookings()))
	}
}

func TestCreateBookingTheWeek(t *testing.T) {
	h := newTestService(t, friday).routes()
	if rec := send(h, http.MethodPost, "/bookings", "application/json", hollis); rec.Code != http.StatusCreated {
		t.Fatalf("POST /bookings with Hollis's booking = %d %q, want 201", rec.Code, rec.Body.String())
	}
	// The rest of the week's Furrow-A and Furrow-B passes, in the order Hollis's program books them.
	week := []string{"furrow-b/SQT/15862", "furrow-a/CWR/15879", "furrow-b/YST/15890", "furrow-a/SQT/15910", "furrow-b/SQT/15923"}
	byID := map[string]PassView{}
	for _, p := range catalog() {
		byID[p.ID] = p
	}
	for i, id := range week {
		wantID := fmt.Sprintf("bk-%d", 4502+i)
		rec := send(h, http.MethodPost, "/bookings", "application/json", bookingBody(byID[id]))
		var got Booking
		if err := json.Unmarshal(rec.Body.Bytes(), &got); rec.Code != http.StatusCreated || err != nil || got.ID != wantID {
			t.Errorf("POST /bookings for %s = %d %q, want 201 and %s", id, rec.Code, rec.Body.String(), wantID)
			continue
		}
		if loc := rec.Header().Get("Location"); loc != "/bookings/"+wantID {
			t.Errorf("POST /bookings for %s: Location = %q, want %q", id, loc, "/bookings/"+wantID)
		}
	}
}

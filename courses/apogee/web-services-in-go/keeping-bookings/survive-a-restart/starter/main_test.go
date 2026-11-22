package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// pass returns one of Furrowlight's passes over Salar Quieto in January 2028.
func pass(id, sat string, day, hour, minute, second int, lastsFor time.Duration, maxElevDeg float64) PassView {
	aos := time.Date(2028, 1, day, hour, minute, second, 0, time.UTC)
	return PassView{ID: id, Pass: Pass{SatID: sat, Station: "SQT", AOS: aos, LOS: aos.Add(lastsFor), MaxElevDeg: maxElevDeg}}
}

// january holds Furrowlight's first passes of January, by AOS.
var january = []PassView{
	pass("furrow-a/SQT/16290", "furrow-a", 3, 14, 11, 55, 10*time.Minute+4*time.Second, 70.5),
	pass("furrow-b/SQT/16288", "furrow-b", 3, 14, 59, 32, 9*time.Minute+33*time.Second, 30.9),
	pass("furrow-a/SQT/16305", "furrow-a", 4, 13, 52, 29, 9*time.Minute+36*time.Second, 33.5),
	pass("furrow-b/SQT/16303", "furrow-b", 4, 14, 39, 33, 10*time.Minute+3*time.Second, 64.7),
}

// body returns the JSON body that books the catalog pass with this ID.
func body(passID string) string {
	for _, p := range january {
		if p.ID == passID {
			return fmt.Sprintf(`{"sat_id":%q,"station":%q,"aos":%q,"los":%q}`,
				p.SatID, p.Station, p.AOS.Format(time.RFC3339), p.LOS.Format(time.RFC3339))
		}
	}
	panic("no pass " + passID + " in the test's catalog")
}

// friday is the service's clock: Friday 17 December 2027, the day of the last deploy of the year.
func friday() time.Time { return time.Date(2027, 12, 17, 10, 0, 0, 0, time.UTC) }

// start runs a new process of the service from cfg, as a deploy does, and returns its server.
func start(t *testing.T, cfg Config) *httptest.Server {
	t.Helper()
	s, err := newService(cfg)
	if err != nil {
		t.Fatalf("newService(a folder of bookings) error = %v, want nil", err)
	}
	srv := httptest.NewServer(s.handler(log.New(io.Discard, "", 0)))
	t.Cleanup(srv.Close)
	return srv
}

// do sends one request to srv and returns the status and the body.
func do(t *testing.T, srv *httptest.Server, method, path, body string) (int, []byte) {
	t.Helper()
	req, err := http.NewRequest(method, srv.URL+path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("%s %s error = %v", method, path, err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s %s: reading the body: %v", method, path, err)
	}
	return res.StatusCode, bytes.TrimSpace(data)
}

// book books the catalog pass with this ID and checks the booking's ID.
func book(t *testing.T, srv *httptest.Server, passID, wantID string) {
	t.Helper()
	status, data := do(t, srv, http.MethodPost, "/bookings", body(passID))
	var b Booking
	if status != http.StatusCreated || json.Unmarshal(data, &b) != nil || b.ID != wantID {
		t.Errorf("POST /bookings (%s) = %d %s, want 201 and booking %s", passID, status, data, wantID)
	}
}

// listIDs asks srv for every booking and returns their IDs in the order it lists them.
func listIDs(t *testing.T, srv *httptest.Server) string {
	t.Helper()
	status, data := do(t, srv, http.MethodGet, "/bookings", "")
	var list []Booking
	if status != http.StatusOK || json.Unmarshal(data, &list) != nil || list == nil {
		t.Fatalf("GET /bookings = %d %s, want 200 and a JSON list", status, data)
	}
	var ids []string
	for _, b := range list {
		ids = append(ids, b.ID)
	}
	return strings.Join(ids, " ")
}

func TestBookingsSurviveRestarts(t *testing.T) {
	cfg := Config{Dir: t.TempDir(), Catalog: january, Now: friday}

	first := start(t, cfg)
	if got := listIDs(t, first); got != "" {
		t.Errorf("GET /bookings of a new service = %q, want no bookings", got)
	}
	book(t, first, "furrow-b/SQT/16288", "bk-4501")
	book(t, first, "furrow-a/SQT/16305", "bk-4502")
	book(t, first, "furrow-a/SQT/16290", "bk-4503")
	first.Close()

	second := start(t, cfg)
	if got, want := listIDs(t, second), "bk-4503 bk-4501 bk-4502"; got != want {
		t.Errorf("GET /bookings after a restart = %q, want %q, by AOS", got, want)
	}
	book(t, second, "furrow-b/SQT/16303", "bk-4504")
	if status, data := do(t, second, http.MethodDelete, "/bookings/bk-4501", ""); status != http.StatusNoContent {
		t.Errorf("DELETE /bookings/bk-4501 = %d %s, want 204", status, data)
	}
	second.Close()

	third := start(t, cfg)
	if got, want := listIDs(t, third), "bk-4503 bk-4502 bk-4504"; got != want {
		t.Errorf("GET /bookings after another restart = %q, want %q, without the canceled bk-4501", got, want)
	}
	if status, _ := do(t, third, http.MethodGet, "/bookings/bk-4501", ""); status != http.StatusNotFound {
		t.Errorf("GET /bookings/bk-4501 after its cancellation and a restart = %d, want 404", status)
	}
}

func TestBrokenBookingsStopTheStart(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bookings.json"), []byte(`[{"id":"bk-4501","sat_id":"furrow-a",`), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := newService(Config{Dir: dir, Catalog: january, Now: friday})
	if err == nil || s != nil {
		t.Errorf("newService(a folder whose bookings.json is cut short) gave a service: %t, error %v; want no service and an error", s != nil, err)
	}
}

package main

import (
	"io"
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

func TestRoutes(t *testing.T) {
	friday := time.Date(2027, 11, 26, 10, 0, 0, 0, time.UTC)
	srv := httptest.NewServer(newTestService(t, friday).routes())
	defer srv.Close()
	requests := []struct {
		method, path string
		status       int
		allow        string // the Allow header of a 405
		body         string // a piece of the body; for HEAD and 204, the body must be empty
	}{
		{"GET", "/health", 200, "", `"status":"ok"`},
		{"DELETE", "/health", 405, "GET, HEAD", "Method Not Allowed"},
		{"GET", "/stations", 200, "", `"name":"Salar Quieto"`},
		{"HEAD", "/stations", 200, "", ""},
		{"POST", "/stations", 405, "GET, HEAD", "Method Not Allowed"},
		{"GET", "/stations/SQT", 200, "", `"code":"SQT"`},
		{"GET", "/stations/ZZZ", 404, "", `no station "ZZZ"`},
		{"POST", "/stations/SQT", 405, "GET, HEAD", "Method Not Allowed"},
		{"GET", "/guide/", 200, "", "Customer guide"},
		{"POST", "/guide/", 405, "GET, HEAD", "Method Not Allowed"},
		{"GET", "/passes", 200, "", `"id":"furrow-b/CWR/15709"`},
		{"DELETE", "/passes", 405, "GET, HEAD", "Method Not Allowed"},
		{"GET", "/passes/furrow-a/SQT/15864", 200, "", `"id":"furrow-a/SQT/15864"`},
		{"GET", "/passes/furrow-a/SQT/15865", 404, "", "no pass furrow-a/SQT/15865"},
		{"PUT", "/passes/furrow-a/SQT/15864", 405, "GET, HEAD", "Method Not Allowed"},
		{"GET", "/bookings/bk-4498", 200, "", `"id":"bk-4498"`},
		{"GET", "/bookings/bk-5499", 404, "", "booking bk-5499: not found"},
		{"DELETE", "/bookings/bk-4497", 204, "", ""},
		{"PUT", "/bookings/bk-4498", 405, "DELETE, GET, HEAD", "Method Not Allowed"},
	}
	for _, rq := range requests {
		t.Run(rq.method+" "+rq.path, func(t *testing.T) {
			req, err := http.NewRequest(rq.method, srv.URL+rq.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			res, err := srv.Client().Do(req)
			if err != nil {
				t.Fatalf("%s %s: %v", rq.method, rq.path, err)
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("%s %s: reading the body: %v", rq.method, rq.path, err)
			}
			if res.StatusCode != rq.status {
				t.Errorf("%s %s = %d %q, want %d", rq.method, rq.path, res.StatusCode, body, rq.status)
			}
			if got := res.Header.Get("Allow"); got != rq.allow {
				t.Errorf("%s %s: Allow = %q, want %q", rq.method, rq.path, got, rq.allow)
			}
			if rq.body == "" && len(body) != 0 {
				t.Errorf("%s %s: body = %q, want none", rq.method, rq.path, body)
			}
			if !strings.Contains(string(body), rq.body) {
				t.Errorf("%s %s: body = %q, want it to hold %q", rq.method, rq.path, body, rq.body)
			}
		})
	}
}

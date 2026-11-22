package main

import (
	"encoding/json"
	"errors"
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

func TestWriteError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		want   ErrorBody
	}{
		{"a 400 from the decoder",
			&APIError{Status: 400, Code: "bad_request", Message: "request body: empty"},
			400, ErrorBody{Code: "bad_request", Message: "request body: empty"}},
		{"a 422, wrapped once",
			fmt.Errorf("booking for furrowlight: %w", &APIError{Status: 422, Code: "invalid_booking", Message: "los must be after aos"}),
			422, ErrorBody{Code: "invalid_booking", Message: "los must be after aos"}},
		{"a 404, wrapped twice",
			fmt.Errorf("get: %w", fmt.Errorf("store: %w", &APIError{Status: 404, Code: "not_found", Message: "booking bk-5499: not found"})),
			404, ErrorBody{Code: "not_found", Message: "booking bk-5499: not found"}},
		{"a bug in the service",
			errors.New("open /srv/bookingd/bookings.json: permission denied"),
			500, ErrorBody{Code: "internal", Message: "internal error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeError(rec, httptest.NewRequest(http.MethodGet, "/bookings/bk-5499", nil), tt.err)
			var got ErrorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != tt.status || got != tt.want {
				t.Errorf("writeError(%v) = %d %q, want %d %+v", tt.err, rec.Code, rec.Body.String(), tt.status, tt.want)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("writeError(%v): Content-Type = %q, want %q", tt.err, ct, "application/json")
			}
			if tt.status == 500 && strings.Contains(rec.Body.String(), "permission denied") {
				t.Errorf("writeError(%v) = %q, which shows the error's own text to the caller", tt.err, rec.Body.String())
			}
		})
	}
}

func TestErrorsAsJSON(t *testing.T) {
	h := newTestService(t, time.Date(2027, 11, 29, 10, 0, 0, 0, time.UTC)).routes()
	tests := []struct {
		method, target, body string
		status               int
		want                 ErrorBody
	}{
		{"GET", "/stations/ZZZ", "", 404, ErrorBody{Code: "not_found", Message: `no station "ZZZ"`}},
		{"GET", "/stations?band=x", "", 400, ErrorBody{Code: "bad_request", Message: `unknown band "x"`}},
		{"GET", "/passes?sat=thistle", "", 400, ErrorBody{Code: "bad_request", Message: `unknown satellite "thistle"`}},
		{"GET", "/passes/furrow-a/SQT/15865", "", 404, ErrorBody{Code: "not_found", Message: "no pass furrow-a/SQT/15865"}},
		{"GET", "/bookings/bk-5499", "", 404, ErrorBody{Code: "not_found", Message: "booking bk-5499: not found"}},
		{"POST", "/bookings", `{"sat_id":`, 400, ErrorBody{Code: "bad_request", Message: "request body: unexpected EOF"}},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.target, func(t *testing.T) {
			rec := send(h, tt.method, tt.target, "application/json", tt.body)
			var got ErrorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != tt.status || got != tt.want {
				t.Errorf("%s %s = %d %q, want %d %+v", tt.method, tt.target, rec.Code, rec.Body.String(), tt.status, tt.want)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("%s %s: Content-Type = %q, want %q", tt.method, tt.target, ct, "application/json")
			}
		})
	}
}

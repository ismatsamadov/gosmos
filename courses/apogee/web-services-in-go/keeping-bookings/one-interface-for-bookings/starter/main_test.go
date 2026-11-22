package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// january holds Furrowlight's first passes of January over Salar Quieto.
var january = []PassView{
	{ID: "furrow-a/SQT/16290", Pass: Pass{SatID: "furrow-a", Station: "SQT",
		AOS: time.Date(2028, 1, 3, 14, 11, 55, 0, time.UTC), LOS: time.Date(2028, 1, 3, 14, 21, 59, 0, time.UTC), MaxElevDeg: 70.5}},
	{ID: "furrow-b/SQT/16288", Pass: Pass{SatID: "furrow-b", Station: "SQT",
		AOS: time.Date(2028, 1, 3, 14, 59, 32, 0, time.UTC), LOS: time.Date(2028, 1, 3, 15, 9, 5, 0, time.UTC), MaxElevDeg: 30.9}},
}

// monday is the service's clock: Monday 13 December 2027, after the deploy.
func monday() time.Time { return time.Date(2027, 12, 13, 10, 0, 0, 0, time.UTC) }

// furrowA is Furrowlight's booking of Furrow-A's first pass of January.
const furrowA = `{"sat_id":"furrow-a","station":"SQT","aos":"2028-01-03T14:11:55Z","los":"2028-01-03T14:21:59Z"}`

// errDisk is what the failing store's disk says; no caller may ever see it.
var errDisk = errors.New("read bookings.json: input/output error")

// failingStore is a store whose disk has failed: every method returns errDisk.
type failingStore struct{}

// Create returns errDisk.
func (f *failingStore) Create(ctx context.Context, b Booking) (Booking, error) {
	return Booking{}, errDisk
}

// Get returns errDisk.
func (f *failingStore) Get(ctx context.Context, id string) (Booking, error) {
	return Booking{}, errDisk
}

// List returns errDisk.
func (f *failingStore) List(ctx context.Context) ([]Booking, error) {
	return nil, errDisk
}

// Delete returns errDisk.
func (f *failingStore) Delete(ctx context.Context, id string) error {
	return errDisk
}

// call sends one request through h and returns the status and the body it got.
func call(h http.Handler, method, path, body string) (int, []byte) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, bytes.TrimSpace(rec.Body.Bytes())
}

func TestHandlersWithMemStore(t *testing.T) {
	s := &Service{Catalog: january, Store: NewMemStore(), Now: monday}
	h := s.handler(log.New(io.Discard, "", 0))
	steps := []struct {
		name, method, path, body string
		status                   int
	}{
		{"Furrowlight books Furrow-A's pass", http.MethodPost, "/bookings", furrowA, http.StatusCreated},
		{"and reads it back", http.MethodGet, "/bookings/bk-4501", "", http.StatusOK},
		{"then cancels it", http.MethodDelete, "/bookings/bk-4501", "", http.StatusNoContent},
		{"and it is gone", http.MethodGet, "/bookings/bk-4501", "", http.StatusNotFound},
	}
	for _, st := range steps {
		t.Run(st.name, func(t *testing.T) {
			status, body := call(h, st.method, st.path, st.body)
			if status != st.status {
				t.Fatalf("%s %s = %d %s, want %d", st.method, st.path, status, body, st.status)
			}
			if status == http.StatusOK || status == http.StatusCreated {
				var b Booking
				if err := json.Unmarshal(body, &b); err != nil || b.ID != "bk-4501" || b.SatID != "furrow-a" {
					t.Errorf("%s %s body = %s, want booking bk-4501 of furrow-a", st.method, st.path, body)
				}
			}
		})
	}
}

func TestHandlersWithFailingStore(t *testing.T) {
	s := &Service{Catalog: january, Store: &failingStore{}, Now: monday}
	h := s.handler(log.New(io.Discard, "", 0))
	for _, rq := range []struct{ method, path, body string }{
		{http.MethodGet, "/bookings/bk-4501", ""},
		{http.MethodPost, "/bookings", furrowA},
	} {
		status, body := call(h, rq.method, rq.path, rq.body)
		if status != http.StatusInternalServerError {
			t.Errorf("%s %s with a failing store = %d, want %d", rq.method, rq.path, status, http.StatusInternalServerError)
		}
		var e ErrorBody
		if err := json.Unmarshal(body, &e); err != nil || e.Code != "internal" || e.Message != "internal error" {
			t.Errorf("%s %s with a failing store: body = %s, want error \"internal\" and message \"internal error\"", rq.method, rq.path, body)
		}
		if strings.Contains(string(body), "input/output") {
			t.Errorf("%s %s: body = %s, which shows the store's own error to the caller", rq.method, rq.path, body)
		}
	}
}

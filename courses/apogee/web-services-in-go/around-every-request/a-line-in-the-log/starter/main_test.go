package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// fakeClock is a clock the test moves by hand.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

// Now returns the clock's time.
func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Add moves the clock forward by d.
func (c *fakeClock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// cirrinesBooking is the booking Cirrine's program made on Wednesday: 236 bytes of JSON with its newline.
const cirrinesBooking = `{"id":"bk-4512","sat_id":"cirrine-1","station":"YST","aos":"2027-12-13T08:26:10Z","los":"2027-12-13T08:36:28Z","note":"Occultation run 3 of 5. Track Cirrine-1 from AOS to LOS, as on Monday: we compare the two passes minute by minute."}` + "\n"

// wednesday answers like the service did on Wednesday; its handlers move the clock as their work would.
func wednesday(clock *fakeClock) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"status":"ok"}`+"\n")
	})
	mux.HandleFunc("GET /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"error":"not_found","message":"booking bk-4511: not found"}`+"\n")
	})
	mux.HandleFunc("GET /passes", func(w http.ResponseWriter, r *http.Request) {
		clock.Add(1204*time.Millisecond + 400*time.Microsecond)
		io.WriteString(w, "[]\n")
	})
	mux.HandleFunc("DELETE /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		// writes nothing at all, which the server sends as 200 with an empty body
	})
	mux.HandleFunc("POST /bookings", func(w http.ResponseWriter, r *http.Request) {
		clock.Add(4*time.Millisecond + 300*time.Microsecond)
		w.Header().Set("Location", "/bookings/bk-4512")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, cirrinesBooking)
	})
	return mux
}

func TestLogRequests(t *testing.T) {
	var buf bytes.Buffer
	clock := &fakeClock{now: time.Date(2027, 12, 8, 10, 40, 0, 0, time.UTC)}
	h := requestID(logRequests(log.New(&buf, "", 0), clock.Now)(wednesday(clock)))
	tests := []struct {
		name, method, target string
		want                 string
	}{
		{"Ingrid's monitor, answered with a body alone", http.MethodGet, "/health", "req-000001 GET /health 200 16B 0s"},
		{"Cirrine's canceled booking", http.MethodGet, "/bookings/bk-4511", "req-000002 GET /bookings/bk-4511 404 61B 0s"},
		{"a slow week of passes, its query left out", http.MethodGet, "/passes?sat=cirrine-1&days=7", "req-000003 GET /passes 200 3B 1.204s"},
		{"a path the mux doesn't know", http.MethodGet, "/pases", "req-000004 GET /pases 404 19B 0s"},
		{"a method the path doesn't take", http.MethodPost, "/health", "req-000005 POST /health 405 19B 0s"},
		{"a cancellation, answered with no body", http.MethodDelete, "/bookings/bk-4503", "req-000006 DELETE /bookings/bk-4503 204 0B 0s"},
		{"Cirrine's new booking", http.MethodPost, "/bookings", "req-000007 POST /bookings 201 236B 4ms"},
		{"a ping answered by writing nothing", http.MethodGet, "/ping", "req-000008 GET /ping 200 0B 0s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(tt.method, tt.target, nil))
			if got, want := buf.String(), tt.want+"\n"; got != want {
				t.Errorf("%s %s logged %q, want %q", tt.method, tt.target, got, want)
			}
		})
	}
}

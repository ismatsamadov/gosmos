package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// cirrineBookings answers like the service did for Cirrine's booking, and counts its calls.
func cirrineBookings(calls *int) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		*calls++
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"id":"bk-4511","sat_id":"cirrine-1","station":"YST","aos":"2027-12-10T08:02:06Z","los":"2027-12-10T08:12:20Z"}`+"\n")
	})
	return mux
}

func TestNoStore(t *testing.T) {
	tests := []struct {
		name, method, path string
		status             int
	}{
		{"Cirrine's booking", http.MethodGet, "/bookings/bk-4511", http.StatusOK},
		{"a path nobody registered", http.MethodGet, "/bokings/bk-4511", http.StatusNotFound},
		{"a method the path doesn't take", http.MethodPost, "/bookings/bk-4511", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			rec := httptest.NewRecorder()
			noStore(cirrineBookings(&calls)).ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))
			res := rec.Result() // the headers as they were when the answer was written
			if res.StatusCode != tt.status {
				t.Errorf("%s %s = %d, want %d", tt.method, tt.path, res.StatusCode, tt.status)
			}
			if got := res.Header.Get("Cache-Control"); got != "no-store" {
				t.Errorf("%s %s: Cache-Control = %q, want %q", tt.method, tt.path, got, "no-store")
			}
		})
	}
}

func TestNoStoreCallsNextOnce(t *testing.T) {
	calls := 0
	h := noStore(cirrineBookings(&calls))
	for want := 1; want <= 3; want++ {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/bookings/bk-4511", nil))
		if calls != want {
			t.Fatalf("after %d requests through noStore, the handler ran %d times, want %d", want, calls, want)
		}
	}
}

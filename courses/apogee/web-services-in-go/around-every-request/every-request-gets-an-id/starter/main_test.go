package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// missingBooking answers every request as the service answers a booking that doesn't exist,
// and keeps the request ID that it finds in each request's context.
func missingBooking(seen *[]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*seen = append(*seen, requestIDFrom(r.Context()))
		writeError(w, r, &APIError{Status: http.StatusNotFound, Code: "not_found", Message: "booking bk-5499: not found"})
	})
}

func TestRequestID(t *testing.T) {
	var seen []string
	h := requestID(missingBooking(&seen))
	for i, want := range []string{"req-000001", "req-000002"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/bookings/bk-5499", nil))
		if got := rec.Result().Header.Get("X-Request-Id"); got != want {
			t.Errorf("request %d: X-Request-Id = %q, want %q", i+1, got, want)
		}
		if len(seen) == i+1 && seen[i] != want {
			t.Errorf("request %d: requestIDFrom(the handler's context) = %q, want %q", i+1, seen[i], want)
		}
		var body ErrorBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("request %d: error body %q doesn't decode: %v", i+1, rec.Body.String(), err)
		}
		if body.RequestID != want {
			t.Errorf("request %d: request_id in the error body = %q, want %q", i+1, body.RequestID, want)
		}
	}
}

func TestRequestIDFrom(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{"a context that never went through requestID", context.Background()},
		{"someone else's value under the plain string key", context.WithValue(context.Background(), "request_id", "req-424242")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestIDFrom(tt.ctx); got != "" {
				t.Errorf("requestIDFrom(ctx) = %q, want %q", got, "")
			}
		})
	}
}

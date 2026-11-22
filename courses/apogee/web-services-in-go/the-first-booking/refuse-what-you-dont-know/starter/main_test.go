package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// describe says what decodeJSON returned: nil, or an *APIError's status, code and message.
func describe(err error) string {
	if err == nil {
		return "nil"
	}
	if e, ok := errors.AsType[*APIError](err); ok {
		return fmt.Sprintf("%d %s %q", e.Status, e.Code, e.Message)
	}
	return fmt.Sprintf("a plain error %q", err)
}

// hollis is the body of Furrowlight's first booking: Furrow-A over Salar Quieto on Monday 6 December.
const hollis = `{"sat_id":"furrow-a","station":"SQT","aos":"2027-12-06T14:02:10Z","los":"2027-12-06T14:11:52Z","note":"PRIORITY"}`

func TestDecodeJSON(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string // the 400's message, or "" when the body is accepted
	}{
		{"Hollis's booking", hollis, ""},
		{"a newline and spaces after it", hollis + "\n  \n", ""},
		{"sat instead of sat_id", strings.Replace(hollis, `"sat_id"`, `"sat"`, 1), `request body: json: unknown field "sat"`},
		{"two bookings in one body", hollis + hollis, "request body: unexpected data after the JSON value"},
		{"a stray brace after it", hollis + "}", "request body: unexpected data after the JSON value"},
		{"an empty body", "", "request body: empty"},
		{"white space alone", " \n\t", "request body: empty"},
		{"a body cut short", `{"sat_id":"furrow-a",`, "request body: unexpected EOF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Booking
			req := httptest.NewRequest(http.MethodPost, "/bookings", strings.NewReader(tt.body))
			err := decodeJSON(httptest.NewRecorder(), req, &b)
			if tt.want == "" {
				if err != nil || b.SatID != "furrow-a" {
					t.Errorf("decodeJSON(%s) = %s with sat_id %q, want nil with %q", tt.name, describe(err), b.SatID, "furrow-a")
				}
				return
			}
			if e, ok := errors.AsType[*APIError](err); !ok || e.Status != 400 || e.Code != "bad_request" || e.Message != tt.want {
				t.Errorf("decodeJSON(%s) = %s, want 400 bad_request %q", tt.name, describe(err), tt.want)
			}
		})
	}
}

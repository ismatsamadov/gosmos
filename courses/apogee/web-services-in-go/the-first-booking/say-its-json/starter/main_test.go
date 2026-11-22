package main

import (
	"errors"
	"fmt"
	"io"
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

// countingReader counts the bytes read through it.
type countingReader struct {
	r io.Reader
	n int64
}

// Read reads from the reader it wraps and counts what it got.
func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}

func TestDecodeJSONMediaType(t *testing.T) {
	const brinewatch = `{"sat_id":"brine-3","station":"LAN","aos":"2027-12-06T09:14:20Z","los":"2027-12-06T09:22:02Z"}`
	tests := []struct {
		name        string
		contentType string // "" leaves the header out
		body        string
		accepted    bool
	}{
		{"JSON", "application/json", brinewatch, true},
		{"JSON with its charset", "application/json; charset=utf-8", brinewatch, true},
		{"JSON in capitals", "Application/JSON", brinewatch, true},
		{"Brinewatch's form", "application/x-www-form-urlencoded", "sat_id=brine-3&station=LAN", false},
		{"JSON sent as plain text", "text/plain; charset=utf-8", brinewatch, false},
		{"JSONP, which only starts the same", "application/jsonp", brinewatch, false},
		{"no Content-Type at all", "", brinewatch, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &countingReader{r: strings.NewReader(tt.body)}
			req := httptest.NewRequest(http.MethodPost, "/bookings", body)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			var b Booking
			err := decodeJSON(httptest.NewRecorder(), req, &b)
			if tt.accepted {
				if err != nil || b.SatID != "brine-3" {
					t.Errorf("decodeJSON(Content-Type %q) = %s with sat_id %q, want nil with %q", tt.contentType, describe(err), b.SatID, "brine-3")
				}
				return
			}
			if e, ok := errors.AsType[*APIError](err); !ok || e.Status != 415 || e.Code != "unsupported_media_type" || e.Message != "send the body as application/json" {
				t.Errorf("decodeJSON(Content-Type %q) = %s, want 415 unsupported_media_type %q", tt.contentType, describe(err), "send the body as application/json")
			}
			if body.n != 0 {
				t.Errorf("decodeJSON(Content-Type %q) read %d bytes of a body it refuses, want 0", tt.contentType, body.n)
			}
		})
	}
}

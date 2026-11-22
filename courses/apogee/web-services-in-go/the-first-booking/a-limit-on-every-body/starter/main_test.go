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

// letters is an endless stream of the letter A, a note that never ends.
type letters struct{}

// Read fills p with the letter A.
func (letters) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'A'
	}
	return len(p), nil
}

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

// withNote returns Hollis's booking with a note of n letters.
func withNote(n int) string {
	return strings.Replace(hollis, `"PRIORITY"`, `"`+strings.Repeat("A", n)+`"`, 1)
}

func TestDecodeJSONLimit(t *testing.T) {
	padded := hollis + strings.Repeat(" ", maxBodyBytes-len(hollis))
	tests := []struct {
		name   string
		body   string
		status int // 0 when the body is accepted
		want   string
	}{
		{"Hollis's booking padded to exactly 65,536 bytes", padded, 0, ""},
		{"one byte more", padded + " ", 413, "request body over 65536 bytes"},
		{"a 70,000-letter note", withNote(70000), 413, "request body over 65536 bytes"},
		{"70,000 spaces after the booking", hollis + strings.Repeat(" ", 70000), 413, "request body over 65536 bytes"},
		{"two bookings, still refused", hollis + hollis, 400, "request body: unexpected data after the JSON value"},
		{"an empty body, still refused", "", 400, "request body: empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Booking
			req := httptest.NewRequest(http.MethodPost, "/bookings", strings.NewReader(tt.body))
			err := decodeJSON(httptest.NewRecorder(), req, &b)
			if tt.status == 0 {
				if err != nil {
					t.Errorf("decodeJSON(%s) = %s, want nil", tt.name, describe(err))
				}
				return
			}
			code := map[int]string{400: "bad_request", 413: "too_large"}[tt.status]
			if e, ok := errors.AsType[*APIError](err); !ok || e.Status != tt.status || e.Code != code || e.Message != tt.want {
				t.Errorf("decodeJSON(%s) = %s, want %d %s %q", tt.name, describe(err), tt.status, code, tt.want)
			}
		})
	}
}

func TestDecodeJSONStopsReading(t *testing.T) {
	// Hollis's booking up to its note's first quote, then a 40 MB note.
	start := strings.TrimSuffix(withNote(0), `"}`)
	note := io.MultiReader(strings.NewReader(start), io.LimitReader(letters{}, 40<<20), strings.NewReader(`"}`))
	body := &countingReader{r: note}
	req := httptest.NewRequest(http.MethodPost, "/bookings", body)
	var b Booking
	err := decodeJSON(httptest.NewRecorder(), req, &b)
	if e, ok := errors.AsType[*APIError](err); !ok || e.Status != 413 {
		t.Errorf("decodeJSON(a 40 MB note) = %s, want 413 too_large", describe(err))
	}
	if body.n > maxBodyBytes+1 {
		t.Errorf("decodeJSON(a 40 MB note) read %d bytes of the body, want at most %d", body.n, maxBodyBytes+1)
	}
}

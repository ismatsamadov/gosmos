package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

// landed is the event for Furrow-A's image from Monday's pass over Salar Quieto.
var landed = Event{
	ID:        "evt-000001",
	Type:      "data.landed",
	BookingID: "bk-4718",
	SatID:     "furrow-a",
	Station:   "SQT",
	File:      "furrow-a-20280117-1416.tif",
	Bytes:     10875000000,
	LandedAt:  time.Date(2028, 1, 17, 14, 27, 52, 0, time.UTC),
}

// The receiver's clock reads 14:27:54 on Monday; the event was signed a second earlier.
const (
	now      = int64(1831732074)
	sentAt   = now - 1
	secret   = "furrowlight webhook secret"
	limitLen = 64 << 10
)

// signed returns a request carrying body to Furrowlight's receiver, signed with key at ts as the service signs it.
func signed(body string, ts int64, key string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/hooks/apogee", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Apogee-Timestamp", strconv.FormatInt(ts, 10))
	r.Header.Set("Apogee-Signature", signature([]byte(key), ts, []byte(body)))
	return r
}

func TestVerifyEvent(t *testing.T) {
	data, _ := json.Marshal(landed)
	body := string(data)
	tests := []struct {
		name string
		req  func() *http.Request
		ok   bool
	}{
		{"Monday's event, a second old", func() *http.Request { return signed(body, sentAt, secret) }, true},
		{"the body changed on the way", func() *http.Request {
			r := signed(body, sentAt, secret)
			r.Body = io.NopCloser(strings.NewReader(strings.Replace(body, "bk-4718", "bk-4719", 1)))
			return r
		}, false},
		{"the timestamp changed on the way", func() *http.Request {
			r := signed(body, sentAt, secret)
			r.Header.Set("Apogee-Timestamp", strconv.FormatInt(sentAt+60, 10))
			return r
		}, false},
		{"signed exactly 300 seconds ago", func() *http.Request { return signed(body, now-300, secret) }, true},
		{"signed 301 seconds ago", func() *http.Request { return signed(body, now-301, secret) }, false},
		{"signed exactly 300 seconds ahead", func() *http.Request { return signed(body, now+300, secret) }, true},
		{"signed 301 seconds ahead", func() *http.Request { return signed(body, now+301, secret) }, false},
		{"signed with Cirrine's secret", func() *http.Request { return signed(body, sentAt, "cirrine webhook secret") }, false},
		{"no timestamp", func() *http.Request {
			r := signed(body, sentAt, secret)
			r.Header.Del("Apogee-Timestamp")
			return r
		}, false},
		{"no signature", func() *http.Request {
			r := signed(body, sentAt, secret)
			r.Header.Del("Apogee-Signature")
			return r
		}, false},
		{"a signature that isn't hex", func() *http.Request {
			r := signed(body, sentAt, secret)
			r.Header.Set("Apogee-Signature", strings.Repeat("zz", 32))
			return r
		}, false},
		{"exactly 64 KiB, padded with spaces", func() *http.Request {
			return signed(body+strings.Repeat(" ", limitLen-len(body)), sentAt, secret)
		}, true},
		{"one byte over 64 KiB", func() *http.Request {
			return signed(body+strings.Repeat(" ", limitLen+1-len(body)), sentAt, secret)
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyEvent(tt.req(), []byte(secret), time.Unix(now, 0))
			switch {
			case tt.ok && err != nil:
				t.Errorf("verifyEvent(%s) = %v, want the event", tt.name, err)
			case tt.ok && got != landed:
				t.Errorf("verifyEvent(%s) = %+v, want %+v", tt.name, got, landed)
			case !tt.ok && err == nil:
				t.Errorf("verifyEvent(%s) = nil error, want it refused", tt.name)
			case !tt.ok && got != (Event{}):
				t.Errorf("verifyEvent(%s) = %+v with its error, want a zero Event: nothing in a refused body is trusted", tt.name, got)
			}
		})
	}
}

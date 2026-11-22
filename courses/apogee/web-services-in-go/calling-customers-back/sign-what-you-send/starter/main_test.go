package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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

// furrowlightSecret is Furrowlight's test webhook secret.
const furrowlightSecret = "furrowlight webhook secret"

// hexMAC is the test's own HMAC-SHA256 of message under key, in lowercase hex.
func hexMAC(key, message string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestSignature(t *testing.T) {
	body, _ := json.Marshal(landed)
	tests := []struct {
		name   string
		secret string
		ts     int64
		body   string
	}{
		{"Monday's event at 14:27:53", furrowlightSecret, 1831732073, string(body)},
		{"the same event a second later", furrowlightSecret, 1831732074, string(body)},
		{"the same event signed for Cirrine", "cirrine webhook secret", 1831732073, string(body)},
		{"an empty body", furrowlightSecret, 1831732073, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := hexMAC(tt.secret, strconv.FormatInt(tt.ts, 10)+"."+tt.body)
			if got := signature([]byte(tt.secret), tt.ts, []byte(tt.body)); got != want {
				t.Errorf("signature(%q, %d, body) = %q, want %q", tt.secret, tt.ts, got, want)
			}
		})
	}
	first := signature([]byte(furrowlightSecret), 1831732073, body)
	later := signature([]byte(furrowlightSecret), 1831732074, body)
	if first == later {
		t.Errorf("signature(%d, body) = signature(%d, body) = %q, want them different: the timestamp must be signed too", 1831732073, 1831732074, first)
	}
}

func TestPostSigns(t *testing.T) {
	type arrival struct{ timestamp, signature, contentType, body string }
	seen := make(chan arrival, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		seen <- arrival{r.Header.Get("Apogee-Timestamp"), r.Header.Get("Apogee-Signature"), r.Header.Get("Content-Type"), string(body)}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	if err := post(context.Background(), srv.Client(), srv.URL, []byte(furrowlightSecret), 1831732073, landed); err != nil {
		t.Fatalf("post(evt-000001) = %v, want nil", err)
	}
	var got arrival
	select {
	case got = <-seen:
	default:
		t.Fatal("post(evt-000001) returned nil, but the receiver got no request")
	}
	if got.timestamp != "1831732073" {
		t.Errorf("post(evt-000001) sent Apogee-Timestamp %q, want %q", got.timestamp, "1831732073")
	}
	if want := hexMAC(furrowlightSecret, "1831732073."+got.body); got.signature != want {
		t.Errorf("post(evt-000001) sent Apogee-Signature %q, want %q, the receiver's own computation over the bytes it got", got.signature, want)
	}
	if got.contentType != "application/json" {
		t.Errorf("post(evt-000001) sent Content-Type %q, want %q", got.contentType, "application/json")
	}
}

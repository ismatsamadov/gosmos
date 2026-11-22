package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"testing"
	"time"
)

// fileDigest is the digest dlinkd gives for Sunday's file.
const fileDigest = "sha256:57f31a45f7f99d89febc9b563983e20056006094293fb00e8d7a03ef78224ce9"

// wantReceipt is the receipt Sunday's landing gets, byte for byte as the customer guide shows it.
const wantReceipt = `{"v":1,"key_id":"r-2028-08","booking_id":"bk-5274","sat_id":"furrow-a","station":"SQT",` +
	`"file":"furrow-a-20280827-1402.tif","bytes":10874220000,"digest":"` + fileDigest + `",` +
	`"landed_at":"2028-08-27T14:13:05Z","signed_at":"2028-09-01T09:00:00Z"}`

// testNotary returns bookingd's notary for the test receipt key r-2028-08, whose seed is the
// SHA-256 of a label, readable on purpose.
func testNotary(t *testing.T) *Notary {
	t.Helper()
	seed := sha256.Sum256([]byte("receipts r-2028-08"))
	n, err := newNotary("r-2028-08", seed[:])
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// sundayBooking is bk-5274: Furrow-A over Salar Quieto on Sunday 27 August.
func sundayBooking() Booking {
	return Booking{
		ID: "bk-5274", SatID: "furrow-a", Station: "SQT",
		AOS: time.Date(2028, 8, 27, 14, 2, 10, 0, time.UTC),
		LOS: time.Date(2028, 8, 27, 14, 11, 52, 0, time.UTC),
	}
}

// sundayLanding is dlinkd's landing for bk-5274's file.
func sundayLanding() Landing {
	return Landing{BookingID: "bk-5274", File: "furrow-a-20280827-1402.tif", Bytes: 10874220000,
		At: time.Date(2028, 8, 27, 14, 13, 5, 0, time.UTC)}
}

// friday is bookingd's clock, stopped at 09:00:00.731 on Friday 1 September.
func friday() time.Time {
	return time.Date(2028, 9, 1, 9, 0, 0, 731_000_000, time.UTC)
}

// short returns s, or its first 8 characters and an ellipsis when it is longer than 50.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

func TestIssue(t *testing.T) {
	n := testNotary(t)
	env, err := issue(n, sundayLanding(), sundayBooking(), fileDigest, friday)
	if err != nil {
		t.Fatalf("issue(Sunday's landing) = %v, want nil", err)
	}
	r, err := verifyEnvelope(map[string]ed25519.PublicKey{"r-2028-08": n.Published()}, env)
	if err != nil {
		t.Fatalf("verifyEnvelope(issue(Sunday's landing)) = %v, want nil", err)
	}
	fields := []struct {
		name      string
		got, want any
	}{
		{"v", r.Version, 1},
		{"key_id", r.KeyID, "r-2028-08"},
		{"booking_id", r.BookingID, "bk-5274"},
		{"sat_id", r.SatID, "furrow-a"},
		{"station", r.Station, "SQT"},
		{"file", r.File, "furrow-a-20280827-1402.tif"},
		{"bytes", r.Bytes, int64(10874220000)},
		{"digest", short(r.Digest), short(fileDigest)},
		{"landed_at", r.LandedAt.Format(time.RFC3339Nano), "2028-08-27T14:13:05Z"},
		{"signed_at", r.SignedAt.Format(time.RFC3339Nano), "2028-09-01T09:00:00Z"},
	}
	for _, f := range fields {
		if f.got != f.want {
			t.Errorf("issue(Sunday's landing): %s = %v, want %v", f.name, f.got, f.want)
		}
	}
}

// TestExactBytes compares the signed receipt with the guide's: field order, no spaces, and
// times to the second.
func TestExactBytes(t *testing.T) {
	env, err := issue(testNotary(t), sundayLanding(), sundayBooking(), fileDigest, friday)
	if err != nil {
		t.Fatalf("issue(Sunday's landing) = %v, want nil", err)
	}
	var e Envelope
	if err := json.Unmarshal(env, &e); err != nil {
		t.Fatalf("issue(Sunday's landing) = an envelope that doesn't decode: %v", err)
	}
	got := string(e.Receipt)
	if got != wantReceipt {
		i := 0
		for i < len(got) && i < len(wantReceipt) && got[i] == wantReceipt[i] {
			i++
		}
		t.Errorf("issue(Sunday's landing) signs %q… from byte %d, want %q…", got[i:min(i+30, len(got))], i, wantReceipt[i:min(i+30, len(wantReceipt))])
	}
}

func TestRefusals(t *testing.T) {
	other := sundayLanding()
	other.BookingID = "bk-5275"
	tests := []struct {
		name    string
		landing Landing
		digest  string
	}{
		{"a digest in capitals", sundayLanding(), "SHA256:" + fileDigest[len("sha256:"):]},
		{"a digest without its hash", sundayLanding(), fileDigest[len("sha256:"):]},
		{"a digest of 63 hex digits", sundayLanding(), fileDigest[:len(fileDigest)-1]},
		{"a landing for bk-5275", other, fileDigest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if env, err := issue(testNotary(t), tt.landing, sundayBooking(), tt.digest, friday); err == nil || env != nil {
				t.Errorf("issue(%s) = %d bytes, %v, want nothing signed and an error", tt.name, len(env), err)
			}
		})
	}
}

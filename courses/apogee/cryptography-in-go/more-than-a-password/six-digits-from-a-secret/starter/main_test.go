package main

import (
	"strings"
	"testing"
	"time"
)

// rfcKey is the secret the test vectors of RFC 4226 and RFC 6238 are made with: 20 ASCII bytes.
const rfcKey = "12345678901234567890"

// mirelasPhone is the test secret in Mirela's authenticator app: 20 ASCII bytes.
const mirelasPhone = "mirela's phone, 2028"

// TestRFC4226 checks the ten codes RFC 4226 lists for counters 0 to 9, in its Appendix D.
func TestRFC4226(t *testing.T) {
	want := []string{"755224", "287082", "359152", "969429", "338314", "254676", "287922", "162583", "399871", "520489"}
	for counter, code := range want {
		if got := hotp([]byte(rfcKey), uint64(counter)); got != code {
			t.Errorf("hotp(%q, %d) = %q, want %q", rfcKey, counter, got, code)
		}
	}
}

// TestRFC6238 checks totp against RFC 6238's SHA-1 vectors, in its Appendix B. The RFC prints
// eight digits; six digits are the same number's last six.
func TestRFC6238(t *testing.T) {
	tests := []struct {
		unix int64
		rfc  string
	}{
		{59, "94287082"},
		{1111111109, "07081804"},
		{1111111111, "14050471"},
		{1234567890, "89005924"},
		{2000000000, "69279037"},
		{20000000000, "65353130"},
	}
	for _, tt := range tests {
		at := time.Unix(tt.unix, 0).UTC()
		if got, want := totp([]byte(rfcKey), at), tt.rfc[2:]; got != want {
			t.Errorf("totp(%q, %s) = %q, want %q, the RFC's %s cut to six digits", rfcKey, at.Format(time.RFC3339), got, want, tt.rfc)
		}
	}
}

// TestMirelasPhone checks the code Mirela's phone shows on Monday 30 October 2028 at 10:00 UTC,
// and those of the steps on either side.
func TestMirelasPhone(t *testing.T) {
	at := time.Date(2028, 10, 30, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		at   time.Time
		want string
	}{
		{at.Add(-30 * time.Second), "462810"},
		{at, "687056"},
		{at.Add(29 * time.Second), "687056"},
		{at.Add(30 * time.Second), "173810"},
	}
	for _, tt := range tests {
		if got := totp([]byte(mirelasPhone), tt.at); got != tt.want {
			t.Errorf("totp(%q, %s) = %q, want %q", mirelasPhone, tt.at.Format(time.RFC3339), got, tt.want)
		}
	}
}

// TestSixDigits checks a thousand codes in a row: six digits each, leading zeros kept.
func TestSixDigits(t *testing.T) {
	for counter := uint64(0); counter < 1000; counter++ {
		code := hotp([]byte(mirelasPhone), counter)
		if len(code) != 6 || strings.Trim(code, "0123456789") != "" {
			t.Fatalf("hotp(%q, %d) = %q, want six digits", mirelasPhone, counter, code)
		}
	}
}

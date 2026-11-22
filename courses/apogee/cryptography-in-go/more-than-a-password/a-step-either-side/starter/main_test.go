package main

import (
	"testing"
	"time"
)

// The test secrets in the authenticator apps: 20 ASCII bytes each, and the RFCs' own.
const (
	ingridsPhone = "ingrid's phone, 2028"
	mirelasPhone = "mirela's phone, 2028"
	rfcKey       = "12345678901234567890"
)

// clock is the console's clock in the tests: a set instant, moved by hand.
type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// testCodes returns a code checker for ingrid, mirela and rfc, whose clock reads at.
func testCodes(at time.Time) *Codes {
	c := &clock{t: at}
	return &Codes{
		secrets: map[string][]byte{"ingrid": []byte(ingridsPhone), "mirela": []byte(mirelasPhone), "rfc": []byte(rfcKey)},
		now:     c.now,
	}
}

// at returns a time on Tuesday 31 October 2028, UTC.
func at(hour, min, sec int) time.Time {
	return time.Date(2028, 10, 31, hour, min, sec, 0, time.UTC)
}

// TestIngridsPhone plays Ingrid's phone, whose clock runs 20 seconds behind the console's, and
// a phone 20 seconds ahead.
func TestIngridsPhone(t *testing.T) {
	behind := totp([]byte(ingridsPhone), at(7, 30, 50))
	if !testCodes(at(7, 31, 10)).Check("ingrid", behind) {
		t.Errorf("Check(%q, %q) at 07:31:10 = false, want true: it is the code a phone 20 seconds behind shows", "ingrid", behind)
	}
	ahead := totp([]byte(ingridsPhone), at(7, 31, 40))
	if !testCodes(at(7, 31, 20)).Check("ingrid", ahead) {
		t.Errorf("Check(%q, %q) at 07:31:20 = false, want true: it is the code a phone 20 seconds ahead shows", "ingrid", ahead)
	}
}

// TestWindow checks the code Ingrid's phone shows at 07:31:05, the step from 07:31:00 to
// 07:31:30, at console times around it: one step either side, and no more.
func TestWindow(t *testing.T) {
	code := totp([]byte(ingridsPhone), at(7, 31, 5))
	tests := []struct {
		name string
		at   time.Time
		want bool
	}{
		{"07:30:29, two steps before", at(7, 30, 29), false},
		{"07:30:30, the step before", at(7, 30, 30), true},
		{"07:31:05, its own step", at(7, 31, 5), true},
		{"07:31:59, the step after", at(7, 31, 59), true},
		{"07:32:00, two steps after", at(7, 32, 0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testCodes(tt.at).Check("ingrid", code); got != tt.want {
				t.Errorf("Check(%q, %q) at %s = %v, want %v", "ingrid", code, tt.at.Format("15:04:05"), got, tt.want)
			}
		})
	}
}

// TestRFC6238 checks the RFC's own codes, cut to six digits, at the RFC's own times.
func TestRFC6238(t *testing.T) {
	tests := []struct {
		unix int64
		code string
	}{
		{59, "287082"},
		{1111111109, "081804"},
		{1111111111, "050471"},
		{1234567890, "005924"},
		{2000000000, "279037"},
		{20000000000, "353130"},
	}
	for _, tt := range tests {
		now := time.Unix(tt.unix, 0).UTC()
		if !testCodes(now).Check("rfc", tt.code) {
			t.Errorf("Check(%q, %q) at %s = false, want true", "rfc", tt.code, now.Format(time.RFC3339))
		}
	}
}

// TestRefused checks codes that must never pass: another account's, malformed ones, and any
// code at all for an account without a secret.
func TestRefused(t *testing.T) {
	now := at(7, 31, 10)
	step := timeStep(now)
	tests := []struct {
		name, user, code string
	}{
		{"mirela's code", "ingrid", totp([]byte(mirelasPhone), now)},
		{"no code", "ingrid", ""},
		{"seven digits", "ingrid", totp([]byte(ingridsPhone), now) + "0"},
		{"the code an empty secret makes, the step before", "kiosk", hotp(nil, step-1)},
		{"the code an empty secret makes", "kiosk", hotp(nil, step)},
		{"the code an empty secret makes, the step after", "kiosk", hotp(nil, step+1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testCodes(now).Check(tt.user, tt.code) {
				t.Errorf("Check(%q, %q) at 07:31:10 = true, want false", tt.user, tt.code)
			}
		})
	}
}

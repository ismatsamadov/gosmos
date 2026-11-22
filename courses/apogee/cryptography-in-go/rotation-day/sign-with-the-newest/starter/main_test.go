package main

import (
	"testing"
	"time"
)

// at returns the instant of an RFC 3339 time.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// The two versions in uplinkd's ring since Monday 9 October, with the season's test secrets.
var (
	k07 = KeyVersion{
		Key:      Key{ID: "k-2028-07", Secret: []byte("uplink test master key, 32 bytes")},
		Activate: at("2028-08-11T00:00:00Z"),
		Retire:   at("2028-10-24T00:00:00Z"),
	}
	k10 = KeyVersion{
		Key:      Key{ID: "k-2028-10", Secret: []byte("uplink next master key, 32 bytes")},
		Activate: at("2028-10-17T00:00:00Z"),
	}
)

func TestSigningKey(t *testing.T) {
	gone := k10
	gone.Revoked = true
	tests := []struct {
		name string
		ring []KeyVersion // in the key file's order; the test also tries it reversed
		now  string
		want string // the version's ID, or the error's text
	}{
		{"a second before Rotation Day", []KeyVersion{k07, k10}, "2028-10-16T23:59:59Z", "k-2028-07"},
		{"Rotation Day, 00:00:00", []KeyVersion{k07, k10}, "2028-10-17T00:00:00Z", "k-2028-10"},
		{"k-2028-10 revoked, 20 October", []KeyVersion{k07, gone}, "2028-10-20T00:00:00Z", "k-2028-07"},
		{"k-2028-10 revoked, the instant k-2028-07 retires", []KeyVersion{k07, gone}, "2028-10-24T00:00:00Z", "no key to sign with at 2028-10-24T00:00:00Z"},
		{"k-2028-10 revoked, 25 October", []KeyVersion{k07, gone}, "2028-10-25T00:00:00Z", "no key to sign with at 2028-10-25T00:00:00Z"},
		{"the night before k-2028-07 activated", []KeyVersion{k07, k10}, "2028-08-10T23:59:59Z", "no key to sign with at 2028-08-10T23:59:59Z"},
	}
	for _, tt := range tests {
		reversed := []KeyVersion{tt.ring[1], tt.ring[0]}
		for _, order := range []struct {
			name string
			ring []KeyVersion
		}{{"key file's order", tt.ring}, {"reversed", reversed}} {
			t.Run(tt.name+", "+order.name, func(t *testing.T) {
				v, err := signingKey(order.ring, at(tt.now))
				got := v.ID
				if err != nil {
					got = err.Error()
				}
				if got != tt.want {
					t.Errorf("signingKey(ring, %s) = %q, want %q", tt.now, got, tt.want)
				}
			})
		}
	}
}

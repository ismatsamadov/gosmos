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

// revoked returns v, revoked.
func revoked(v KeyVersion) KeyVersion {
	v.Revoked = true
	return v
}

// TestRehearsal replays Mirela's rehearsal: a command sealed with k-2028-07 at 23:59:58 reaches
// Cirrine-1 at 00:00:03 on Rotation Day, three seconds after k-2028-10 took over.
func TestRehearsal(t *testing.T) {
	arrives := at("2028-10-17T00:00:03Z")
	if !accepts([]KeyVersion{k07, k10}, "k-2028-07", arrives) {
		t.Errorf("accepts(ring, %q, %s) = false, want true: the pass would be lost", "k-2028-07", arrives.Format(time.RFC3339))
	}
}

func TestAccepts(t *testing.T) {
	ring := []KeyVersion{k07, k10}
	tests := []struct {
		name string
		ring []KeyVersion
		id   string
		now  string
		want bool
	}{
		{"k-2028-10, held a week before it activates", ring, "k-2028-10", "2028-10-09T12:00:00Z", true},
		{"k-2028-07, one second before it retires", ring, "k-2028-07", "2028-10-23T23:59:59Z", true},
		{"k-2028-07 at its retire instant", ring, "k-2028-07", "2028-10-24T00:00:00Z", false},
		{"k-2028-10, which has no retire time", ring, "k-2028-10", "2029-03-01T00:00:00Z", true},
		{"k-2027-03, which the ring doesn't hold", ring, "k-2027-03", "2028-10-17T00:00:03Z", false},
		{"k-2028-10, revoked", []KeyVersion{k07, revoked(k10)}, "k-2028-10", "2028-10-20T00:00:00Z", false},
		{"k-2028-07, revoked before it retires", []KeyVersion{revoked(k07), k10}, "k-2028-07", "2028-10-10T00:00:00Z", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := accepts(tt.ring, tt.id, at(tt.now)); got != tt.want {
				t.Errorf("accepts(ring, %q, %s) = %v, want %v", tt.id, tt.now, got, tt.want)
			}
		})
	}
}

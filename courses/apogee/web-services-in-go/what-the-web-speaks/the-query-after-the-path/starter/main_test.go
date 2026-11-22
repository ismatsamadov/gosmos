package main

import (
	"net/url"
	"testing"
)

func TestPassesPath(t *testing.T) {
	tests := []struct {
		name               string
		sat, station, from string
		want               string
	}{
		{"Furrow-A over Salar Quieto", "furrow-a", "SQT", "", "/passes?sat=furrow-a&station=SQT"},
		{"every pass over Salar Quieto", "", "SQT", "", "/passes?station=SQT"},
		{"every pass there is", "", "", "", "/passes"},
		{"Hollis's 09:00, keys sorted", "furrow-a", "", "2027-12-06T09:00:00+01:00", "/passes?from=2027-12-06T09%3A00%3A00%2B01%3A00&sat=furrow-a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := passesPath(tt.sat, tt.station, tt.from)
			if got != tt.want {
				t.Errorf("passesPath(%q, %q, %q) = %q, want %q", tt.sat, tt.station, tt.from, got, tt.want)
			}
		})
	}
}

func TestPassesPathReadsBack(t *testing.T) {
	const from = "2027-12-06T09:00:00+01:00"
	path := passesPath("furrow-a", "SQT", from)
	u, err := url.Parse(path)
	if err != nil {
		t.Fatalf("url.Parse(%q) error = %v, want nil", path, err)
	}
	q := u.Query()
	for _, kv := range [][2]string{{"sat", "furrow-a"}, {"station", "SQT"}, {"from", from}} {
		if got := q.Get(kv[0]); got != kv[1] {
			t.Errorf("url.Parse(%q).Query().Get(%q) = %q, want %q", path, kv[0], got, kv[1])
		}
	}
}

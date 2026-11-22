package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// saved is bookings.json as the service left it before a restart: four of Furrowlight's
// January bookings, in the order the store happened to write them.
const saved = `[
{"id":"bk-4503","sat_id":"furrow-a","station":"SQT","aos":"2028-01-04T13:52:29Z","los":"2028-01-04T14:02:05Z"},
{"id":"bk-4501","sat_id":"furrow-a","station":"SQT","aos":"2028-01-03T14:11:55Z","los":"2028-01-03T14:21:59Z"},
{"id":"bk-4504","sat_id":"furrow-b","station":"SQT","aos":"2028-01-04T14:39:33Z","los":"2028-01-04T14:49:36Z","note":"PRIORITY"},
{"id":"bk-4502","sat_id":"furrow-b","station":"SQT","aos":"2028-01-03T14:59:32Z","los":"2028-01-03T15:09:05Z"}
]
`

// furrowA5 is a booking of Furrow-A's pass over Salar Quieto on 5 January, which no file holds yet.
var furrowA5 = Booking{SatID: "furrow-a", Station: "SQT",
	AOS: time.Date(2028, 1, 5, 13, 33, 24, 0, time.UTC), LOS: time.Date(2028, 1, 5, 13, 41, 56, 0, time.UTC)}

// writeSaved writes text as bookings.json in a new folder and returns the file's path.
func writeSaved(t *testing.T, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bookings.json")
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// nextID opens the store at path and returns the ID its next booking gets.
func nextID(t *testing.T, path string) string {
	t.Helper()
	s, err := OpenFileStore(path)
	if err != nil {
		t.Fatalf("OpenFileStore(bookings.json) error = %v, want nil", err)
	}
	b, err := s.Create(context.Background(), furrowA5)
	if err != nil {
		t.Fatalf("Create(furrow-a on 5 January) error = %v, want nil", err)
	}
	return b.ID
}

func TestOpenAMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bookings.json")
	if got := nextID(t, path); got != "bk-4501" {
		t.Errorf("first booking of a store with no file = %q, want %q", got, "bk-4501")
	}
}

func TestOpenRestoresBookings(t *testing.T) {
	path := writeSaved(t, saved)
	s, err := OpenFileStore(path)
	if err != nil {
		t.Fatalf("OpenFileStore(four saved bookings) error = %v, want nil", err)
	}
	ctx := context.Background()
	list, err := s.List(ctx)
	if err != nil || len(list) != 4 {
		t.Fatalf("List() after opening four saved bookings = %d bookings, %v, want 4, nil", len(list), err)
	}
	for i, want := range []string{"bk-4501", "bk-4502", "bk-4503", "bk-4504"} {
		if list[i].ID != want {
			t.Errorf("List()[%d].ID = %q, want %q", i, list[i].ID, want)
		}
	}
	if b, err := s.Get(ctx, "bk-4504"); err != nil || b.Note != "PRIORITY" {
		t.Errorf("Get(%q) = %+v, %v, want Furrow-B's booking with note PRIORITY", "bk-4504", b, err)
	}
	made, err := s.Create(ctx, furrowA5)
	if err != nil || made.ID != "bk-4505" {
		t.Fatalf("next Create = %q, %v, want %q, nil", made.ID, err, "bk-4505")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var after []Booking
	if err := json.Unmarshal(data, &after); err != nil || len(after) != 5 {
		t.Errorf("after one more booking, bookings.json holds %d bookings, want 5: the store must not save over what it opened", len(after))
	}
}

func TestOpenContinuesAboveYusrasIDs(t *testing.T) {
	path := writeSaved(t, `[{"id":"bk-4491","sat_id":"cirrine-1","station":"YST","aos":"2027-11-29T08:09:40Z","los":"2027-11-29T08:19:55Z"}]`)
	if got := nextID(t, path); got != "bk-4501" {
		t.Errorf("next ID after Yusra's bk-4491 = %q, want %q: the service's own count starts at bk-4501", got, "bk-4501")
	}
}

func TestOpenRefusesBrokenFiles(t *testing.T) {
	tests := []struct{ name, text string }{
		{"a file cut short", `[{"id":"bk-4501","sat_id":"furrow-a",`},
		{"a field the store doesn't know", `[{"id":"bk-4501","satellite":"furrow-a"}]`},
		{"two arrays", "[]\n[]\n"},
		{"one booking, not an array", `{"id":"bk-4501","sat_id":"furrow-a"}`},
		{"one ID saved twice", `[{"id":"bk-4501","sat_id":"furrow-a"},{"id":"bk-4501","sat_id":"furrow-b"}]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := OpenFileStore(writeSaved(t, tt.text))
			if err == nil || s != nil {
				t.Errorf("OpenFileStore(%s) gave a store: %t, error %v; want no store and an error", tt.name, s != nil, err)
			}
		})
	}
}

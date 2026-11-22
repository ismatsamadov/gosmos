package main

import (
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
	"time"
)

// recorderLog is three lines of the van's booking log.
const recorderLog = `{"id":"bk-2071","sat_id":"furrow-a","station":"SQT","aos":"2027-08-24T14:02:10Z","los":"2027-08-24T14:11:52Z","note":"priority"}
{"id":"bk-2072","sat_id":"furrow-b","station":"SQT","aos":"2027-08-24T14:49:41Z","los":"2027-08-24T14:59:13Z"}
{"id":"bk-2073","sat_id":"brine-3","station":"LAN","aos":"2027-08-25T02:17:05Z","los":"2027-08-25T02:26:31Z"}
`

// bookingIDs lists the IDs of bookings, in order.
func bookingIDs(bookings []Booking) []string {
	var ids []string
	for _, b := range bookings {
		ids = append(ids, b.ID)
	}
	return ids
}

func TestDecodeStream(t *testing.T) {
	all := []string{"bk-2071", "bk-2072", "bk-2073"}
	for _, tc := range []struct {
		name, log string
		want      []string
	}{
		{"three lines", recorderLog, all},
		{"blank lines and tabs", "\n\t" + strings.ReplaceAll(recorderLog, "\n", "\n\n\t  "), all},
		{"no final newline", strings.TrimSuffix(recorderLog, "\n"), all},
		{"empty log", "", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeStream(strings.NewReader(tc.log))
			if ids := bookingIDs(got); !slices.Equal(ids, tc.want) || err != nil {
				t.Errorf("decodeStream(%s) = %q, %v, want %q, nil", tc.name, ids, err, tc.want)
			}
		})
	}
	got, err := decodeStream(strings.NewReader(recorderLog))
	if len(got) != 3 || err != nil {
		t.Fatalf("decodeStream(log) = %d bookings, %v, want 3, nil", len(got), err)
	}
	want := Booking{ID: "bk-2071", SatID: "furrow-a", Station: "SQT", Note: "priority",
		AOS: time.Date(2027, 8, 24, 14, 2, 10, 0, time.UTC), LOS: time.Date(2027, 8, 24, 14, 11, 52, 0, time.UTC)}
	b := got[0]
	if b.ID != want.ID || b.SatID != want.SatID || b.Station != want.Station || b.Note != want.Note || !b.AOS.Equal(want.AOS) || !b.LOS.Equal(want.LOS) {
		t.Errorf("decodeStream(log)[0] = %+v, want %+v", b, want)
	}
}

func TestDecodeStreamStopsAtDamage(t *testing.T) {
	for _, tc := range []struct {
		name, log string
		want      []string
		err       error // nil means any error
	}{
		{"line cut off", recorderLog + `{"id":"bk-2074","sat_id":"fur`, []string{"bk-2071", "bk-2072", "bk-2073"}, io.ErrUnexpectedEOF},
		{"garbled second line", `{"id":"bk-2071"}` + "\n" + `{"id":"bk-2072",,}` + "\n", []string{"bk-2071"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeStream(strings.NewReader(tc.log))
			switch {
			case err == nil:
				t.Errorf("decodeStream(%s) error = nil, want an error", tc.name)
			case tc.err != nil && !errors.Is(err, tc.err):
				t.Errorf("decodeStream(%s) error = %v, want io.ErrUnexpectedEOF", tc.name, err)
			}
			if ids := bookingIDs(got); !slices.Equal(ids, tc.want) {
				t.Errorf("decodeStream(%s) = %q, want the complete bookings %q", tc.name, ids, tc.want)
			}
		})
	}
}

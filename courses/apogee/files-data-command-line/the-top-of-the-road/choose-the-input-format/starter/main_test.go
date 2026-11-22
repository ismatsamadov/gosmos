package main

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

// export is Yusra's export: a JSON array of bookings.
const export = `[
{"id":"bk-3410","sat_id":"furrow-a","station":"SQT","aos":"2027-10-21T14:02:10Z","los":"2027-10-21T14:13:52Z","note":"PRIORITY, FIELD 7"},
{"id":"bk-3413","sat_id":"brine-3","station":"SQT","aos":"2027-10-21T15:40:05Z","los":"2027-10-21T15:51:31Z"}
]`

// sheet is Hollis's sheet: the same bookings as CSV, its columns in his order.
const sheet = `STATION,ID,SATELLITE,AOS,LOS,NOTE
SQT,bk-3410,furrow-a,2027-10-21T14:02:10Z,2027-10-21T14:13:52Z,"PRIORITY, FIELD 7"
SQT,bk-3413,brine-3,2027-10-21T15:40:05Z,2027-10-21T15:51:31Z,
`

// sameBookings reports whether a and b hold the same bookings, times compared as instants.
func sameBookings(a, b []Booking) bool {
	return slices.EqualFunc(a, b, func(x, y Booking) bool {
		return x.ID == y.ID && x.SatID == y.SatID && x.Station == y.Station &&
			x.AOS.Equal(y.AOS) && x.LOS.Equal(y.LOS) && x.Note == y.Note
	})
}

func TestReadBookings(t *testing.T) {
	want := []Booking{
		{ID: "bk-3410", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2027, 10, 21, 14, 2, 10, 0, time.UTC), LOS: time.Date(2027, 10, 21, 14, 13, 52, 0, time.UTC), Note: "PRIORITY, FIELD 7"},
		{ID: "bk-3413", SatID: "brine-3", Station: "SQT", AOS: time.Date(2027, 10, 21, 15, 40, 5, 0, time.UTC), LOS: time.Date(2027, 10, 21, 15, 51, 31, 0, time.UTC)},
	}
	for _, tc := range []struct{ format, input string }{{"json", export}, {"csv", sheet}} {
		t.Run(tc.format, func(t *testing.T) {
			got, err := readBookings(strings.NewReader(tc.input), tc.format)
			if err != nil || !sameBookings(got, want) {
				t.Errorf("readBookings(r, %q) = %v, %v, want %v, nil", tc.format, got, err, want)
			}
		})
	}
}

func TestReadBookingsRefuses(t *testing.T) {
	for _, tc := range []struct{ name, format, input, mention string }{
		{"unknown format", "xml", export, "xml"},
		{"no format", "", export, ""},
		{"sheet named as json", "json", sheet, ""},
		{"export named as csv", "csv", export, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readBookings(strings.NewReader(tc.input), tc.format)
			if got != nil || err == nil {
				t.Errorf("readBookings(r, %q) = %v, %v, want nil and an error", tc.format, got, err)
			} else if !strings.Contains(err.Error(), tc.mention) {
				t.Errorf("readBookings(r, %q) error = %q, want it to name %q", tc.format, err, tc.mention)
			}
		})
	}
	problem := errors.New("stick unplugged")
	for _, format := range []string{"json", "csv"} {
		_, err := readBookings(readerFunc(func([]byte) (int, error) { return 0, problem }), format)
		if !errors.Is(err, problem) {
			t.Errorf("readBookings(failing reader, %q) error = %v, want the reader's error", format, err)
		}
	}
}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the test function.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

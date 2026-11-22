package main

import (
	"encoding/json"
	"testing"
)

const (
	bk2041 = `{"id":"bk-2041","sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z"}`
	bk2042 = `{"id":"bk-2042","sat_id":"furrow-b","station":"SQT","aos":"2027-08-11T14:49:40Z","los":"2027-08-11T14:58:31Z"}`
)

func TestDecodeBookings(t *testing.T) {
	got, err := decodeBookings([]byte("[" + bk2041 + ",\n" + bk2042 + "]"))
	if len(got) != 2 || got[0].ID != "bk-2041" || got[1].ID != "bk-2042" || got[1].SatID != "furrow-b" || err != nil {
		t.Errorf("decodeBookings(bk-2041, bk-2042) = %+v, %v, want bk-2041 then bk-2042, nil", got, err)
	}
}

func TestDecodeBookingsEmptyDay(t *testing.T) {
	for _, tc := range []struct{ name, data string }{
		{"an empty array", "[]"},
		{"null", "null"},
		{"an empty array with spaces", " [ ]\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeBookings([]byte(tc.data))
			out, _ := json.Marshal(got)
			if string(out) != "[]" || err != nil {
				t.Errorf("json.Marshal(decodeBookings(%q)) = %s, with error %v, want [], nil", tc.data, out, err)
			}
		})
	}
}

func TestDecodeBookingsInvalid(t *testing.T) {
	for _, tc := range []struct{ name, data string }{
		{"a number for an id", "[" + bk2041 + `,{"id":2042}]`},
		{"one booking, not a list", bk2041},
		{"cut short", "[" + bk2041 + ","},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeBookings([]byte(tc.data))
			if got != nil || err == nil {
				t.Errorf("decodeBookings(%s) = %+v, %v, want nil, an error", tc.name, got, err)
			}
		})
	}
}

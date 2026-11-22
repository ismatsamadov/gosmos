package main

import (
	"encoding/json"
	"testing"
	"time"
)

const export = `{"id":"bk-2042","sat_id":"furrow-b","station":"SQT","aos":"2027-08-11T14:49:40Z","los":"2027-08-11T14:58:31Z"}`

// bk2042 is the booking the export line holds.
var bk2042 = Booking{
	ID:      "bk-2042",
	SatID:   "furrow-b",
	Station: "SQT",
	AOS:     time.Date(2027, 8, 11, 14, 49, 40, 0, time.UTC),
	LOS:     time.Date(2027, 8, 11, 14, 58, 31, 0, time.UTC),
}

// same reports whether two bookings hold the same values.
func same(a, b Booking) bool {
	return a.ID == b.ID && a.SatID == b.SatID && a.Station == b.Station && a.AOS.Equal(b.AOS) && a.LOS.Equal(b.LOS)
}

func TestBookingDecodes(t *testing.T) {
	var got Booking
	err := json.Unmarshal([]byte(export), &got)
	if !same(got, bk2042) || err != nil {
		t.Errorf("json.Unmarshal(export, &b) gives %+v, %v, want %+v, nil", got, err, bk2042)
	}
}

func TestBookingEncodes(t *testing.T) {
	data, err := json.Marshal(bk2042)
	if string(data) != export || err != nil {
		t.Errorf("json.Marshal(bk-2042) = %s, %v, want %s, nil", data, err, export)
	}
	var back Booking
	err = json.Unmarshal(data, &back)
	if !same(back, bk2042) || err != nil {
		t.Errorf("json.Unmarshal(json.Marshal(bk-2042)) gives %+v, %v, want %+v, nil", back, err, bk2042)
	}
}

package main

import (
	"encoding/json"
	"testing"
	"time"
)

// bk2041 is the booking the van confirms.
var bk2041 = Booking{
	ID:      "bk-2041",
	SatID:   "furrow-a",
	Station: "SQT",
	AOS:     time.Date(2027, 8, 11, 14, 2, 10, 0, time.UTC),
	LOS:     time.Date(2027, 8, 11, 14, 11, 52, 0, time.UTC),
	Note:    "priority",
}

func TestEncodeBooking(t *testing.T) {
	want := `{"id":"bk-2041","sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z","note":"priority"}`
	got, err := encodeBooking(bk2041)
	if string(got) != want || err != nil {
		t.Errorf("encodeBooking(bk-2041) = %s, %v, want %s, nil", got, err, want)
	}
}

func TestEncodeBookingEscapes(t *testing.T) {
	quoted := bk2041
	quoted.Note = `HOLLIS: "PRIORITY" MEANS PRIORITY`
	want := `{"id":"bk-2041","sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z","note":"HOLLIS: \"PRIORITY\" MEANS PRIORITY"}`
	got, err := encodeBooking(quoted)
	if string(got) != want || err != nil {
		t.Errorf("encodeBooking(quoted note) = %s, %v, want %s, nil", got, err, want)
	}
	var back Booking
	if err := json.Unmarshal(got, &back); err != nil || back.Note != quoted.Note || back.ID != quoted.ID || !back.AOS.Equal(quoted.AOS) {
		t.Errorf("json.Unmarshal(encodeBooking(quoted note)) gives %+v, %v, want %+v, nil", back, err, quoted)
	}
}

func TestEncodeBookingRefuses(t *testing.T) {
	typo := bk2041
	typo.AOS = time.Date(20270, 8, 11, 14, 2, 10, 0, time.UTC)
	got, err := encodeBooking(typo)
	if got != nil || err == nil {
		t.Errorf("encodeBooking(year 20270) = %s, %v, want nil, an error", got, err)
	}
}

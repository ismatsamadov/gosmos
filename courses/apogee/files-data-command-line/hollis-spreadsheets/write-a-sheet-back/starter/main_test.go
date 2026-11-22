package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"reflect"
	"testing"
)

// correctedRows is the booking sheet as Yusra corrected it.
var correctedRows = [][]string{
	{"ID", "SATELLITE", "STATION", "AOS", "LOS", "NOTE"},
	{"bk-2117", "furrow-a", "SQT", "2027-09-07T13:58:02Z", "2027-09-07T14:07:40Z", "PRIORITY, SPRAYING FROM 06:00"},
	{"bk-2118", "furrow-b", "SQT", "2027-09-07T14:45:31Z", "2027-09-07T14:55:03Z", `MOVED BY "HOLLIS"`},
	{"bk-2119", "furrow-a", "YST", "2027-09-08T05:11:47Z", "2027-09-08T05:21:15Z", "CLOUD OVER THE NORTH FIELDS\nIMAGE ANYWAY"},
	{"bk-2120", "furrow-b", "YST", "2027-09-08T05:59:20Z", "2027-09-08T06:08:51Z", ""},
}

func TestWriteSheet(t *testing.T) {
	want := "ID,SATELLITE,STATION,AOS,LOS,NOTE\n" +
		"bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,\"PRIORITY, SPRAYING FROM 06:00\"\n" +
		"bk-2118,furrow-b,SQT,2027-09-07T14:45:31Z,2027-09-07T14:55:03Z,\"MOVED BY \"\"HOLLIS\"\"\"\n" +
		"bk-2119,furrow-a,YST,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,\"CLOUD OVER THE NORTH FIELDS\nIMAGE ANYWAY\"\n" +
		"bk-2120,furrow-b,YST,2027-09-08T05:59:20Z,2027-09-08T06:08:51Z,\n"
	var out bytes.Buffer
	if err := writeSheet(&out, correctedRows); out.String() != want || err != nil {
		t.Fatalf("writeSheet(corrected sheet) wrote %q, %v, want %q, nil", out.String(), err, want)
	}
	back, err := csv.NewReader(&out).ReadAll()
	if !reflect.DeepEqual(back, correctedRows) || err != nil {
		t.Errorf("reading the written sheet = %q, %v, want %q, nil", back, err, correctedRows)
	}
	out.Reset()
	if err := writeSheet(&out, nil); out.Len() != 0 || err != nil {
		t.Errorf("writeSheet(no rows) wrote %q, %v, want nothing, nil", out.String(), err)
	}
	problem := errors.New("stick full")
	err = writeSheet(writerFunc(func(p []byte) (int, error) { return 0, problem }), correctedRows)
	if !errors.Is(err, problem) {
		t.Errorf("writeSheet(failing writer) error = %v, want %v", err, problem)
	}
}

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the supplied test writer.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

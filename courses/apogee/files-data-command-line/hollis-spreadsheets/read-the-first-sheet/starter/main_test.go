package main

import (
	"encoding/csv"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// septemberSheet is the start of Hollis's booking sheet for September.
const septemberSheet = "ID,SATELLITE,STATION,AOS,LOS,NOTE\n" +
	"bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,PRIORITY\n" +
	"bk-2118,furrow-b,SQT,2027-09-07T14:45:31Z,2027-09-07T14:55:03Z,\n"

func TestReadSheet(t *testing.T) {
	want := [][]string{
		{"ID", "SATELLITE", "STATION", "AOS", "LOS", "NOTE"},
		{"bk-2117", "furrow-a", "SQT", "2027-09-07T13:58:02Z", "2027-09-07T14:07:40Z", "PRIORITY"},
		{"bk-2118", "furrow-b", "SQT", "2027-09-07T14:45:31Z", "2027-09-07T14:55:03Z", ""},
	}
	got, err := readSheet(strings.NewReader(septemberSheet))
	if !reflect.DeepEqual(got, want) || err != nil {
		t.Errorf("readSheet(September) = %q, %v, want %q, nil", got, err, want)
	}
	got, err = readSheet(strings.NewReader(strings.TrimSuffix(septemberSheet, "\n")))
	if !reflect.DeepEqual(got, want) || err != nil {
		t.Errorf("readSheet(no final newline) = %q, %v, want %q, nil", got, err, want)
	}
	got, err = readSheet(strings.NewReader(""))
	if len(got) != 0 || err != nil {
		t.Errorf("readSheet(empty sheet) = %q, %v, want no records, nil", got, err)
	}
}

func TestReadSheetRefusesBadQuotes(t *testing.T) {
	for _, tc := range []struct {
		name, row, wantName string
		want                error
	}{
		{"quote never closed", `bk-2119,furrow-a,YST,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,"PRIORITY` + "\n", "csv.ErrQuote", csv.ErrQuote},
		{"quote inside a cell", `bk-2119,furrow-a,YST,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,CALL "HOLLIS"` + "\n", "csv.ErrBareQuote", csv.ErrBareQuote},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readSheet(strings.NewReader(septemberSheet + tc.row))
			if len(got) != 0 || !errors.Is(err, tc.want) {
				t.Errorf("readSheet(%s) = %d records, %v, want none and %s", tc.name, len(got), err, tc.wantName)
			}
		})
	}
}

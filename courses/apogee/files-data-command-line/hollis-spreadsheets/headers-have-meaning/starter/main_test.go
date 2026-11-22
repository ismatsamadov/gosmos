package main

import (
	"slices"
	"strings"
	"testing"
)

func TestColumnIndex(t *testing.T) {
	for _, tc := range []struct {
		name   string
		header []string
		want   int
	}{
		{"usual order", []string{"ID", "SATELLITE", "STATION", "AOS", "LOS", "NOTE"}, 2},
		{"moved to the front", []string{"STATION", "ID", "SATELLITE", "AOS", "LOS", "NOTE"}, 0},
		{"moved to the end", []string{"ID", "SATELLITE", "AOS", "LOS", "NOTE", "STATION"}, 5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := columnIndex(tc.header, "STATION")
			if got != tc.want || err != nil {
				t.Errorf("columnIndex(%q, STATION) = %d, %v, want %d, nil", tc.header, got, err, tc.want)
			}
		})
	}
	for _, tc := range []struct {
		name   string
		header []string
	}{
		{"missing", []string{"ID", "SATELLITE", "AOS", "LOS", "NOTE"}},
		{"twice", []string{"ID", "STATION", "SATELLITE", "AOS", "LOS", "STATION"}},
		{"empty header", []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := columnIndex(tc.header, "STATION")
			if err == nil || !strings.Contains(err.Error(), "STATION") {
				t.Errorf("columnIndex(%q, STATION) error = %v, want an error naming STATION", tc.header, err)
			}
		})
	}
}

func TestStationColumn(t *testing.T) {
	want := []string{"SQT", "YST"}
	for _, tc := range []struct{ name, sheet string }{
		{"usual order", "ID,SATELLITE,STATION,AOS,LOS,NOTE\n" +
			"bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,PRIORITY\n" +
			"bk-2119,furrow-a,YST,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,\n"},
		{"STATION first", "STATION,ID,SATELLITE,AOS,LOS,NOTE\n" +
			"SQT,bk-2117,furrow-a,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,PRIORITY\n" +
			"YST,bk-2119,furrow-a,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,\n"},
	} {
		got, err := stationColumn(strings.NewReader(tc.sheet))
		if !slices.Equal(got, want) || err != nil {
			t.Errorf("stationColumn(%s) = %q, %v, want %q, nil", tc.name, got, err, want)
		}
	}
}

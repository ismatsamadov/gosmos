package main

import (
	"slices"
	"strings"
	"testing"
)

// header is the first line of Hollis's booking sheet.
const header = "ID,SATELLITE,STATION,AOS,LOS,NOTE\n"

func TestReadNotes(t *testing.T) {
	for _, tc := range []struct{ name, row, note string }{
		{"comma", `bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,"PRIORITY, SPRAYING FROM 06:00"`,
			"PRIORITY, SPRAYING FROM 06:00"},
		{"doubled quotes", `bk-2118,furrow-b,SQT,2027-09-07T14:45:31Z,2027-09-07T14:55:03Z,"HOLLIS SAYS ""URGENT"""`,
			`HOLLIS SAYS "URGENT"`},
		{"line break", "bk-2119,furrow-a,YST,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,\"CLOUD OVER THE NORTH FIELDS\nIMAGE ANYWAY\"",
			"CLOUD OVER THE NORTH FIELDS\nIMAGE ANYWAY"},
		{"no note", "bk-2120,furrow-b,YST,2027-09-08T05:59:20Z,2027-09-08T06:08:51Z,", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readNotes(strings.NewReader(header + tc.row + "\n"))
			if want := []string{tc.note}; !slices.Equal(got, want) || err != nil {
				t.Errorf("readNotes(%s) = %q, %v, want %q, nil", tc.name, got, err, want)
			}
		})
	}
	sheet := header +
		`bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,"PRIORITY, SPRAYING FROM 06:00"` + "\n" +
		"bk-2119,furrow-a,YST,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,\"CLOUD OVER THE NORTH FIELDS\nIMAGE ANYWAY\"\n" +
		"bk-2120,furrow-b,YST,2027-09-08T05:59:20Z,2027-09-08T06:08:51Z,ROUTINE\n"
	want := []string{"PRIORITY, SPRAYING FROM 06:00", "CLOUD OVER THE NORTH FIELDS\nIMAGE ANYWAY", "ROUTINE"}
	got, err := readNotes(strings.NewReader(sheet))
	if !slices.Equal(got, want) || err != nil {
		t.Errorf("readNotes(three bookings) = %q, %v, want %q, nil", got, err, want)
	}
	if got, err := readNotes(strings.NewReader("")); len(got) != 0 || err != nil {
		t.Errorf("readNotes(empty sheet) = %q, %v, want no notes, nil", got, err)
	}
	if got, err := readNotes(strings.NewReader(header + `bk-2121,furrow-a,CWR,2027-09-09T07:02:12Z,2027-09-09T07:11:40Z,"PRIORITY` + "\n")); err == nil {
		t.Errorf("readNotes(quote never closed) = %q, nil, want an error", got)
	}
}

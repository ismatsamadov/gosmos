package main

import (
	"encoding/csv"
	"errors"
	"slices"
	"strings"
	"testing"
)

// header is the first line of Hollis's booking sheet.
const header = "ID,SATELLITE,STATION,AOS,LOS,NOTE\n"

// firstCells lists the first cell of every record.
func firstCells(records [][]string) []string {
	var cells []string
	for _, record := range records {
		cells = append(cells, record[0])
	}
	return cells
}

func TestCheckedRows(t *testing.T) {
	good := header +
		"bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,PRIORITY\n" +
		"bk-2118,furrow-b,SQT,2027-09-07T14:45:31Z,2027-09-07T14:55:03Z,\n"
	got, err := checkedRows(strings.NewReader(good))
	if ids, want := firstCells(got), []string{"ID", "bk-2117", "bk-2118"}; !slices.Equal(ids, want) || err != nil {
		t.Errorf("checkedRows(good sheet) = records %q, %v, want %q, nil", ids, err, want)
	}

	ragged := header +
		"bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z,PRIORITY\n" +
		"bk-2118,furrow-b,SQT,2027-09-07T14:45:31Z,\n" +
		"bk-2119,furrow-a,YST,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z,\n" +
		"bk-2120,furrow-b,YST,2027-09-08T05:59:20Z,2027-09-08T06:08:51Z,CLOUDY, IMAGE ANYWAY\n" +
		"bk-2121,furrow-a,CWR,2027-09-09T07:02:12Z,2027-09-09T07:11:40Z,\n"
	got, err = checkedRows(strings.NewReader(ragged))
	if ids, want := firstCells(got), []string{"ID", "bk-2117", "bk-2119", "bk-2121"}; !slices.Equal(ids, want) {
		t.Errorf("checkedRows(short and long rows) = records %q, want %q", ids, want)
	}
	if err == nil || !errors.Is(err, csv.ErrFieldCount) || !strings.Contains(err.Error(), "line 3") || !strings.Contains(err.Error(), "line 5") {
		t.Errorf("checkedRows(short and long rows) error = %v, want field-count errors naming line 3 and line 5", err)
	}

	five := "ID,SATELLITE,STATION,AOS,LOS\n" +
		"bk-2117,furrow-a,SQT,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z\n"
	got, err = checkedRows(strings.NewReader(five))
	if len(got) != 0 || !errors.Is(err, csv.ErrFieldCount) || !strings.Contains(err.Error(), "line 1") {
		t.Errorf("checkedRows(five columns) = %d records, %v, want none and an error naming line 1", len(got), err)
	}

	broken := header + `bk-2122,furrow-b,CWR,2027-09-09T07:49:30Z,2027-09-09T07:58:57Z,"PRIORITY` + "\n"
	got, err = checkedRows(strings.NewReader(broken))
	if len(got) != 0 || !errors.Is(err, csv.ErrQuote) {
		t.Errorf("checkedRows(quote never closed) = %d records, %v, want none and csv.ErrQuote", len(got), err)
	}
}

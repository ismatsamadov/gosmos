package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// missedPasses is the missing-pass list as Hollis sent it back: notes added and STATION moved.
const missedPasses = "NOTE,STATION,ID,SATELLITE,AOS,LOS\n" +
	"\"NOTHING, NOT EVEN THE START\",SQT,bk-2117,furrow-a,2027-09-07T13:58:02Z,2027-09-07T14:07:40Z\n" +
	"\"LA GRANDE AGAIN?\nSAME AS LAST WEEK\",SQT,bk-2124,furrow-b,2027-09-08T14:31:05Z,2027-09-08T14:40:44Z\n" +
	",YST,bk-2119,furrow-a,2027-09-08T05:11:47Z,2027-09-08T05:21:15Z\n" +
	"\"STOWED? \"\"WIND 74\"\" IN THE LOG\",CWR,bk-2122,furrow-b,2027-09-09T07:49:30Z,2027-09-09T07:58:57Z\n" +
	"SPRAYING DELAYED,SQT,bk-2130,furrow-a,2027-09-09T13:41:12Z,2027-09-09T13:50:47Z\n"

func TestSummarizeSheet(t *testing.T) {
	for _, tc := range []struct{ name, sheet, want string }{
		{"Hollis's sheet", missedPasses, "STATION,MISSED\nCWR,1\nSQT,3\nYST,1\nTOTAL,5\n"},
		{"usual order", "ID,SATELLITE,STATION,AOS,LOS,NOTE\n" +
			"bk-2125,furrow-a,SQT,2027-09-08T15:18:40Z,2027-09-08T15:28:09Z,\n" +
			"bk-2127,furrow-a,LAN,2027-09-08T22:06:51Z,2027-09-08T22:16:20Z,RAIN\n",
			"STATION,MISSED\nLAN,1\nSQT,1\nTOTAL,2\n"},
		{"header only", "NOTE,STATION,ID,SATELLITE,AOS,LOS\n", "STATION,MISSED\nTOTAL,0\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := summarizeSheet(strings.NewReader(tc.sheet), &out)
			if out.String() != tc.want || err != nil {
				t.Errorf("summarizeSheet(%s) wrote %q, %v, want %q, nil", tc.name, out.String(), err, tc.want)
			}
		})
	}
}

func TestSummarizeSheetRefuses(t *testing.T) {
	for _, tc := range []struct{ name, sheet, mention string }{
		{"no STATION column", "NOTE,ID,SATELLITE\nNOTHING,bk-2117,furrow-a\n", "STATION"},
		{"two STATION columns", "STATION,ID,STATION\nSQT,bk-2117,SQT\n", "STATION"},
		{"quote never closed", "NOTE,STATION,ID\n\"NOTHING,SQT,bk-2117\n", ""},
		{"empty sheet", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := summarizeSheet(strings.NewReader(tc.sheet), &out)
			switch {
			case err == nil:
				t.Errorf("summarizeSheet(%s) error = nil, want an error", tc.name)
			case !strings.Contains(err.Error(), tc.mention):
				t.Errorf("summarizeSheet(%s) error = %v, want one naming %q", tc.name, err, tc.mention)
			}
			if out.Len() != 0 {
				t.Errorf("summarizeSheet(%s) wrote %q, want nothing", tc.name, out.String())
			}
		})
	}
	problem := errors.New("stick full")
	err := summarizeSheet(strings.NewReader(missedPasses), writerFunc(func(p []byte) (int, error) { return 0, problem }))
	if !errors.Is(err, problem) {
		t.Errorf("summarizeSheet(failing writer) error = %v, want %v", err, problem)
	}
}

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the supplied test writer.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

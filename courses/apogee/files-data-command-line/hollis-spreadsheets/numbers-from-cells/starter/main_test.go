package main

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// passSheet is Hollis's pass sheet for Tuesday and Wednesday.
const passSheet = "ID,STATION,ELEV,NOTE\n" +
	"bk-2131,SQT,41.7,\n" +
	"bk-2132,SQT,63.5,\"CLEAR SKY,\nGOOD FOR IMAGES\"\n" +
	"bk-2133,YST,8,\n"

func TestReadPassRows(t *testing.T) {
	want := []PassRow{{"bk-2131", 41.7}, {"bk-2132", 63.5}, {"bk-2133", 8}}
	got, err := readPassRows(strings.NewReader(passSheet))
	if !slices.Equal(got, want) || err != nil {
		t.Errorf("readPassRows(pass sheet) = %v, %v, want %v, nil", got, err, want)
	}
	for _, tc := range []struct {
		name, sheet, prefix, wantName string
		want                          error
	}{
		{"a word", passSheet + "bk-2134,CWR,HIGH,\"HOLLIS: \"\"HIGH\"\" MEANS ABOVE 60\"\n", "row 5, ELEV: ", "strconv.ErrSyntax", strconv.ErrSyntax},
		{"after a two-line note", "ID,STATION,ELEV,NOTE\nbk-2132,SQT,63.5,\"CLEAR SKY,\nGOOD FOR IMAGES\"\nbk-2134,CWR,HIGH,\n", "row 3, ELEV: ", "strconv.ErrSyntax", strconv.ErrSyntax},
		{"an empty cell", "ID,STATION,ELEV,NOTE\nbk-2135,LAN,,\n", "row 2, ELEV: ", "strconv.ErrSyntax", strconv.ErrSyntax},
		{"out of range", "ID,STATION,ELEV,NOTE\nbk-2136,PEL,1e400,\n", "row 2, ELEV: ", "strconv.ErrRange", strconv.ErrRange},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readPassRows(strings.NewReader(tc.sheet))
			if len(got) != 0 || err == nil || !strings.HasPrefix(err.Error(), tc.prefix) || !errors.Is(err, tc.want) {
				t.Errorf("readPassRows(%s) = %v, %v, want no rows and an error starting %q that wraps %s", tc.name, got, err, tc.prefix, tc.wantName)
			}
		})
	}
}

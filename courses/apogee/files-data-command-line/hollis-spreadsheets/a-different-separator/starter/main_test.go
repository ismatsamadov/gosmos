package main

import (
	"encoding/csv"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestReadControllerCSV(t *testing.T) {
	header := []string{"TIME", "DISH", "AZ", "EL", "NOTE"}
	for _, tc := range []struct {
		name, export string
		want         [][]string
	}{
		{"semicolons", "TIME;DISH;AZ;EL;NOTE\n2027-09-05T03:14:07Z;LG;183.25;41.70;\n",
			[][]string{header, {"2027-09-05T03:14:07Z", "LG", "183.25", "41.70", ""}}},
		{"commas are text", "TIME;DISH;AZ;EL;NOTE\n2027-09-05T03:14:08Z;LG;183.31;41.66;WIND 41, GUSTS 55\n",
			[][]string{header, {"2027-09-05T03:14:08Z", "LG", "183.31", "41.66", "WIND 41, GUSTS 55"}}},
		{"quoted semicolon", "TIME;DISH;AZ;EL;NOTE\n2027-09-05T03:14:09Z;LG;183.37;41.62;\"CHECKED; NO STOW\"\n",
			[][]string{header, {"2027-09-05T03:14:09Z", "LG", "183.37", "41.62", "CHECKED; NO STOW"}}},
		{"comment lines", "# La Grande pointing table, 2027-09-05\nTIME;DISH;AZ;EL;NOTE\n# wind rising\n2027-09-05T03:14:07Z;LG;183.25;41.70;\n",
			[][]string{header, {"2027-09-05T03:14:07Z", "LG", "183.25", "41.70", ""}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readControllerCSV(strings.NewReader(tc.export))
			if !reflect.DeepEqual(got, tc.want) || err != nil {
				t.Errorf("readControllerCSV(%s) = %q, %v, want %q, nil", tc.name, got, err, tc.want)
			}
		})
	}
	for _, tc := range []struct {
		name, line, wantName string
		want                 error
	}{
		{"quote never closed", "2027-09-05T03:14:10Z;LG;183.43;41.58;\"CHECKED\n", "csv.ErrQuote", csv.ErrQuote},
		{"quote inside a cell", "2027-09-05T03:14:10Z;LG;183.43;41.58;SEE \"MANUAL\"\n", "csv.ErrBareQuote", csv.ErrBareQuote},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readControllerCSV(strings.NewReader("TIME;DISH;AZ;EL;NOTE\n" + tc.line))
			if !errors.Is(err, tc.want) {
				t.Errorf("readControllerCSV(%s) error = %v, want %s", tc.name, err, tc.wantName)
			}
		})
	}
}

package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"
)

const (
	note = "# night shift 2027-07-06, Paz"
	lg   = "2027-07-06T03:14:07Z\tLG\t183.25\t41.70"
	ch   = "2027-07-06T03:14:08Z  CH  183.31   41.66"
	bad  = "2027-07-06T03:14:09Z\tLG\t183.27\t41..70"
)

// wholeLog returns 200,000 lines in shifts of 1,000, with one doubled decimal point on line 123,502.
func wholeLog() string {
	shift := note + "\n" + strings.Repeat(lg+"\n", 998) + "\n"
	broken := note + "\n" + strings.Repeat(lg+"\n", 500) + bad + "\n" + strings.Repeat(lg+"\n", 497) + "\n"
	return strings.Repeat(shift, 123) + broken + strings.Repeat(shift, 76)
}

func TestReadElevations(t *testing.T) {
	got, err := readElevations(strings.NewReader(note + "\n" + lg + "\n\n" + ch + "\n"))
	if len(got) != 2 || math.Abs(got[0]-41.70) > 0.01 || math.Abs(got[1]-41.66) > 0.01 || err != nil {
		t.Errorf("readElevations(good log) = %v, %v, want [41.7 41.66], nil", got, err)
	}
	for _, tc := range []struct {
		name, log, want string
	}{
		{"a doubled decimal point after a note and a blank line", note + "\n" + lg + "\n\n" + bad + "\n" + ch + "\n", `line 4: strconv.ParseFloat: parsing "41..70": invalid syntax`},
		{"a missing column on the first line", "2027-07-06T03:14:10Z\tCH\t183.31\n" + lg + "\n", "line 1: 3 columns instead of 4"},
		{"a doubled decimal point deep in the whole log", wholeLog(), `line 123502: strconv.ParseFloat: parsing "41..70": invalid syntax`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readElevations(strings.NewReader(tc.log))
			if got != nil || err == nil || err.Error() != tc.want {
				t.Errorf("readElevations(%s) = %v, %v, want nil, %s", tc.name, got, err, tc.want)
			}
		})
	}
	_, err = readElevations(strings.NewReader(note + "\n" + bad + "\n"))
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("errors.Is(readElevations(doubled decimal point), strconv.ErrSyntax) = false, want true; the error was %v", err)
	}
}

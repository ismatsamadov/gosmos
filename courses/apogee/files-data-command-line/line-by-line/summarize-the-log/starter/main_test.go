package main

import (
	"errors"
	"io"
	"maps"
	"math"
	"strconv"
	"strings"
	"testing"
)

const note = "# night shift 2027-07-06, Paz\n"

// reading returns one line of the controller log for a dish at an elevation.
func reading(dish, elevDeg string) string {
	return "2027-07-06T03:14:07Z\t" + dish + "\t183.25\t" + elevDeg + "\n"
}

func TestSummarize(t *testing.T) {
	for _, tc := range []struct {
		name, log string
		want      Summary
	}{
		{"an empty log", "", Summary{}},
		{"notes only", note + "\n\t# wind 40 km/h\n", Summary{}},
		{"both dishes", note + reading("LG", "41.70") + reading("CH", "41.66") + "\n" + reading("LG", "52.10") + reading("CH", "38.20") + reading("LG", "47.35"),
			Summary{Readings: 5, ByDish: map[string]int{"LG": 3, "CH": 2}, MaxElevDeg: 52.10}},
		{"the highest first", reading("LG", "63.50") + reading("LG", "12.00"),
			Summary{Readings: 2, ByDish: map[string]int{"LG": 2}, MaxElevDeg: 63.50}},
		{"all below the horizon", reading("LG", "-0.80") + reading("LG", "-0.35") + reading("LG", "-0.60"),
			Summary{Readings: 3, ByDish: map[string]int{"LG": 3}, MaxElevDeg: -0.35}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := summarize(strings.NewReader(tc.log))
			if got.Readings != tc.want.Readings || !maps.Equal(got.ByDish, tc.want.ByDish) || math.Abs(got.MaxElevDeg-tc.want.MaxElevDeg) > 0.01 || err != nil {
				t.Errorf("summarize(%s) = %+v, %v, want %+v, nil", tc.name, got, err, tc.want)
			}
		})
	}
}

func TestSummarizeBadRecords(t *testing.T) {
	for _, tc := range []struct {
		name, log, want string
	}{
		{"a doubled decimal point", note + reading("LG", "41.70") + "\n" + reading("LG", "41..70"), `line 4: strconv.ParseFloat: parsing "41..70": invalid syntax`},
		{"a missing column", reading("CH", "41.66") + "2027-07-06T03:14:10Z\tCH\t183.31\n", "line 2: 3 columns instead of 4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := summarize(strings.NewReader(tc.log))
			if got.Readings != 0 || len(got.ByDish) != 0 || err == nil || err.Error() != tc.want {
				t.Errorf("summarize(%s) = %+v, %v, want zero Summary, %s", tc.name, got, err, tc.want)
			}
		})
	}
	if _, err := summarize(strings.NewReader(reading("LG", "41..70"))); !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("errors.Is(summarize(doubled decimal point), strconv.ErrSyntax) = false, want true; the error was %v", err)
	}
	problem := errors.New("stick unplugged")
	damaged := io.MultiReader(strings.NewReader(reading("LG", "41.70")), readerFunc(func([]byte) (int, error) { return 0, problem }))
	if _, err := summarize(damaged); !errors.Is(err, problem) {
		t.Errorf("summarize(damaged copy) error = %v, want %v", err, problem)
	}
}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the supplied test function.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

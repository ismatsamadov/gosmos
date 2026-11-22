package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestReadReportID(t *testing.T) {
	for _, tc := range []struct {
		name, report string
		want         int64
	}{
		{"image report", `{"version":2,"kind":"image","id":9007199254741003,"payload":{"chunks":2000,"missing":9}}`, 9007199254741003},
		{"its neighbor", `{"version":2,"kind":"station","id":9007199254741005,"payload":{"passes":6,"missing":31}}`, 9007199254741005},
		{"a small id", `{"id":42}`, 42},
		{"the largest int64", `{"id":9223372036854775807}`, 9223372036854775807},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readReportID(strings.NewReader(tc.report))
			if got != tc.want || err != nil {
				t.Errorf("readReportID(%s) = %d, %v, want %d, nil", tc.name, got, err, tc.want)
			}
		})
	}
}

func TestReadReportIDRefuses(t *testing.T) {
	for _, tc := range []struct {
		name, report, wantName string
		want                   error // nil means any error
	}{
		{"a fraction", `{"id":9007199254741003.5}`, "strconv.ErrSyntax", strconv.ErrSyntax},
		{"beyond 64 bits", `{"id":9223372036854775808}`, "strconv.ErrRange", strconv.ErrRange},
		{"no id", `{"version":2,"kind":"image","payload":{}}`, "", nil},
		{"id as text", `{"id":"9007199254741003"}`, "", nil},
		{"cut off", `{"id":90071992`, "", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readReportID(strings.NewReader(tc.report))
			switch {
			case err == nil:
				t.Errorf("readReportID(%s) error = nil, want an error", tc.name)
			case tc.want != nil && !errors.Is(err, tc.want):
				t.Errorf("readReportID(%s) error = %v, want one wrapping %s", tc.name, err, tc.wantName)
			}
		})
	}
}

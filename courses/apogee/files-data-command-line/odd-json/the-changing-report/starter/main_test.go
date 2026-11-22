package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// The reports in the van's log: two from August in version 2, one from July in version 1.
const (
	imageV2   = `{"version":2,"kind":"image","id":9007199254741003,"payload":{"sat_id":"furrow-a","station":"SQT","file":"fa-0824-1402.img","chunks":2000,"missing":9}}`
	stationV2 = `{"version":2,"kind":"station","id":9007199254741005,"payload":{"station":"SQT","day":"2027-08-24","passes":6,"missing":31}}`
	imageV1   = `{"version":1,"id":9007199254740995,"payload":{"sat_id":"furrow-b","station":"SQT","file":"fb-0719-0631.img","chunks":2000,"missing":0}}`
)

// describe writes a report with the value its payload points to.
func describe(r Report) string {
	switch p := r.Payload.(type) {
	case *ImageReport:
		return fmt.Sprintf("{ID:%d Kind:%q Payload:&%+v}", r.ID, r.Kind, *p)
	case *StationReport:
		return fmt.Sprintf("{ID:%d Kind:%q Payload:&%+v}", r.ID, r.Kind, *p)
	}
	return fmt.Sprintf("{ID:%d Kind:%q Payload:%v}", r.ID, r.Kind, r.Payload)
}

func TestNormalizeReport(t *testing.T) {
	for _, tc := range []struct {
		name, report string
		want         Report
	}{
		{"version 2 image", imageV2, Report{ID: 9007199254741003, Kind: "image",
			Payload: &ImageReport{SatID: "furrow-a", Station: "SQT", File: "fa-0824-1402.img", Chunks: 2000, Missing: 9}}},
		{"version 2 station", stationV2 + "\n", Report{ID: 9007199254741005, Kind: "station",
			Payload: &StationReport{Station: "SQT", Day: "2027-08-24", Passes: 6, Missing: 31}}},
		{"version 1 image from July", imageV1, Report{ID: 9007199254740995, Kind: "image",
			Payload: &ImageReport{SatID: "furrow-b", Station: "SQT", File: "fb-0719-0631.img", Chunks: 2000}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeReport(strings.NewReader(tc.report))
			if describe(got) != describe(tc.want) || err != nil {
				t.Errorf("normalizeReport(%s) = %s, %v, want %s, nil", tc.name, describe(got), err, describe(tc.want))
			}
		})
	}
}

func TestNormalizeReportRefuses(t *testing.T) {
	for _, tc := range []struct {
		name, report, prefix string // the error must start with prefix
	}{
		{"version 3", `{"version":3,"kind":"image","id":9007199254741007,"payload":{"sat_id":"furrow-a","parts":[1,2]}}`,
			"report 9007199254741007: unsupported version 3"},
		{"no version", `{"kind":"image","id":9007199254741009,"payload":{"sat_id":"furrow-b"}}`,
			"report 9007199254741009: unsupported version 0"},
		{"unknown kind", `{"version":2,"kind":"billing","id":9007199254741015,"payload":{"amount":1200}}`,
			`report 9007199254741015: unknown report kind "billing"`},
		{"unknown payload field", `{"version":2,"kind":"station","id":9007199254741013,"payload":{"station":"SQT","lost":31}}`,
			"report 9007199254741013: "},
		{"unknown envelope field", `{"version":2,"kind":"image","id":9007199254741017,"sent":"2027-08-26T06:00:00Z","payload":{}}`, ""},
		{"two reports", imageV2 + "\n" + stationV2, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := normalizeReport(strings.NewReader(tc.report))
			switch {
			case err == nil:
				t.Errorf("normalizeReport(%s) error = nil, want an error", tc.name)
			case !strings.HasPrefix(err.Error(), tc.prefix):
				t.Errorf("normalizeReport(%s) error = %v, want one starting %q", tc.name, err, tc.prefix)
			}
		})
	}
	_, err := normalizeReport(strings.NewReader(`{"version":2,"kind":"image","id":9007199254741011,"payload":{"sat_id":"furrow-a","chunks":"2000"}}`))
	var typeErr *json.UnmarshalTypeError
	if err == nil || !strings.HasPrefix(err.Error(), "report 9007199254741011: ") || !errors.As(err, &typeErr) {
		t.Errorf("normalizeReport(chunks as text) error = %v, want one starting %q that wraps the decoding error", err, "report 9007199254741011: ")
	}
}

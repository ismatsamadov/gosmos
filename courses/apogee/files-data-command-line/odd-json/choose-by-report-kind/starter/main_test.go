package main

import (
	"strings"
	"testing"
)

// The reports Furrowlight sent about Tuesday's passes at Salar Quieto.
const (
	imageReport   = `{"version":2,"kind":"image","id":9007199254741003,"payload":{"sat_id":"furrow-a","station":"SQT","file":"fa-0824-1402.img","chunks":2000,"missing":9}}`
	stationReport = `{"version":2,"kind":"station","id":9007199254741005,"payload":{"station":"SQT","day":"2027-08-24","passes":6,"missing":31}}`
)

func TestDecodeReport(t *testing.T) {
	got, err := decodeReport([]byte(imageReport))
	image, ok := got.(*ImageReport)
	if err != nil || !ok {
		t.Fatalf("decodeReport(image report) = %T, %v, want *ImageReport, nil", got, err)
	}
	wantImage := ImageReport{SatID: "furrow-a", Station: "SQT", File: "fa-0824-1402.img", Chunks: 2000, Missing: 9}
	if *image != wantImage {
		t.Errorf("decodeReport(image report) = %+v, want %+v", *image, wantImage)
	}
	got, err = decodeReport([]byte(stationReport))
	station, ok := got.(*StationReport)
	if err != nil || !ok {
		t.Fatalf("decodeReport(station report) = %T, %v, want *StationReport, nil", got, err)
	}
	wantStation := StationReport{Station: "SQT", Day: "2027-08-24", Passes: 6, Missing: 31}
	if *station != wantStation {
		t.Errorf("decodeReport(station report) = %+v, want %+v", *station, wantStation)
	}
}

func TestDecodeReportRefuses(t *testing.T) {
	for _, tc := range []struct{ name, report, mention string }{
		{"unknown kind", `{"version":2,"kind":"billing","id":9007199254741015,"payload":{"amount":1200}}`, "billing"},
		{"chunks as text", `{"version":2,"kind":"image","id":9007199254741011,"payload":{"sat_id":"furrow-a","chunks":"2000"}}`, ""},
		{"station payload as a list", `{"version":2,"kind":"station","id":9007199254741019,"payload":[6,31]}`, ""},
		{"no payload", `{"version":2,"kind":"station","id":9007199254741021}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := decodeReport([]byte(tc.report))
			switch {
			case err == nil:
				t.Errorf("decodeReport(%s) error = nil, want an error", tc.name)
			case !strings.Contains(err.Error(), tc.mention):
				t.Errorf("decodeReport(%s) error = %v, want one naming %q", tc.name, err, tc.mention)
			}
		})
	}
}

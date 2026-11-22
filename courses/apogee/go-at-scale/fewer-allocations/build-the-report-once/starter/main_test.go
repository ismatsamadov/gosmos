package main

import (
	"fmt"
	"strings"
	"testing"
)

// contactsAt returns the contacts of take-over day's draft at the
// station with code, by AOS.
func contactsAt(draft []Contact, code string) []Contact {
	var out []Contact
	for _, c := range draft {
		if c.Station == code {
			out = append(out, c)
		}
	}
	return out
}

// wantReport is a host's report as the tests expect it, line by line.
func wantReport(code string, contacts []Contact) []string {
	lines := []string{fmt.Sprintf("%s contacts: %d", code, len(contacts))}
	for _, c := range contacts {
		lines = append(lines, fmt.Sprintf("%s dish %d", c.ID, c.Dish))
	}
	return append(lines, "")
}

func TestStationReport(t *testing.T) {
	draft := takeOverDraft()
	tests := []struct {
		name     string
		code     string
		contacts []Contact
	}{
		{"Cape Wren's day, one dish", "CWR", contactsAt(draft, "CWR")},
		{"Ystmark's day, two dishes", "YST", contactsAt(draft, "YST")},
		{"Vorrel Island before it goes live", "VRL", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := fmt.Sprintf("stationReport(%q, %d contacts)", tt.code, len(tt.contacts))
			got := strings.Split(stationReport(tt.code, tt.contacts), "\n")
			want := wantReport(tt.code, tt.contacts)
			for i := range min(len(got), len(want)) {
				if got[i] != want[i] {
					t.Fatalf("%s line %d = %q, want %q", call, i+1, got[i], want[i])
				}
			}
			if len(got) != len(want) {
				t.Fatalf("%s has %d lines, want %d, each ending in a newline", call, len(got)-1, len(want)-1)
			}
		})
	}
}

func TestStationReportAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	cwr := contactsAt(takeOverDraft(), "CWR")
	var report string
	allocs := testing.AllocsPerRun(10, func() {
		report = stationReport("CWR", cwr)
	})
	if strings.Count(report, "\n") != len(cwr)+1 {
		return // TestStationReport says what is wrong
	}
	if allocs > 2 {
		t.Errorf("stationReport(\"CWR\", %d contacts) made %.0f allocations, want at most 2: one buffer, grown once for the whole report", len(cwr), allocs)
	}
}

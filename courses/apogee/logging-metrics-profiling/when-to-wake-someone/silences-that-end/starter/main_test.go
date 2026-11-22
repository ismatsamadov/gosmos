package main

import (
	"slices"
	"testing"
	"time"
)

// at returns hh:mm UTC on the given day of January 2029.
func at(day, h, m int) time.Time {
	return time.Date(2029, time.January, day, h, m, 0, 0, time.UTC)
}

// gen is the generator's silence: Lantern Atoll, Thursday 11 January from 09:00 to 11:00,
// while its new part is fitted with the station off the air.
var gen = Silence{Station: "LAN", From: at(11, 9, 0), Until: at(11, 11, 0)}

func TestSilenced(t *testing.T) {
	// Tomás's mute of 28 December, had it ended at 08:00.
	mute := Silence{Station: "LAN", From: time.Date(2028, 12, 28, 7, 28, 0, 0, time.UTC), Until: time.Date(2028, 12, 28, 8, 0, 0, 0, time.UTC)}
	tests := []struct {
		name     string
		what     string // the silences, as a message names them
		silences []Silence
		station  string
		at       time.Time
		want     bool
	}{
		{"08:59, before it", "[gen]", []Silence{gen}, "LAN", at(11, 8, 59), false},
		{"09:00, its start included", "[gen]", []Silence{gen}, "LAN", at(11, 9, 0), true},
		{"10:59, inside it", "[gen]", []Silence{gen}, "LAN", at(11, 10, 59), true},
		{"11:00, its end excluded", "[gen]", []Silence{gen}, "LAN", at(11, 11, 0), false},
		{"Ystmark at 10:00", "[gen]", []Silence{gen}, "YST", at(11, 10, 0), false},
		{"after an older silence that has ended", "[mute gen]", []Silence{mute, gen}, "LAN", at(11, 10, 0), true},
		{"no silences", "nil", nil, "LAN", at(11, 10, 0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := silenced(tt.silences, tt.station, tt.at); got != tt.want {
				t.Errorf("silenced(%s, %q, %s) = %v, want %v", tt.what, tt.station, tt.at.Format("Jan 2 15:04"), got, tt.want)
			}
		})
	}
}

// Ingrid's rules from Monday 8 January.
var (
	passFailures = Rule{Name: "pass-failures", Severity: Critical, FireAt: 2, ClearBelow: 1, Repeat: 24 * time.Hour}
	failedPass   = Rule{Name: "failed-pass", Severity: Warning, FireAt: 1, ClearBelow: 1}
	diskFull     = Rule{Name: "disk-full", Severity: Critical, FireAt: 0.90, ClearBelow: 0.85, For: 15 * time.Minute, Repeat: 12 * time.Hour}
)

// rota is the on-call rota around the generator's silence, oldest first.
var rota = []Shift{
	{From: at(1, 8, 0), Name: "No. 039"},
	{From: at(8, 8, 0), Name: "Nkem Halloran"},
	{From: at(15, 8, 0), Name: "Mirela Dunn"},
}

// sample is a value of a rule's metric at station on Thursday 11 January at hh:mm.
func sample(station string, h, m int, v float64) Sample {
	return Sample{At: at(11, h, m), Station: station, Value: v}
}

// TestSilenceHoldsAlerts tries the generator's silence before Thursday, with the alerts the
// fitting could bring, through the pager.
func TestSilenceHoldsAlerts(t *testing.T) {
	tests := []struct {
		name    string
		rule    Rule
		samples []Sample
		want    []string
	}{
		{"a failed pass at 10:59 is held", failedPass,
			[]Sample{sample("LAN", 10, 59, 1)}, []string{}},
		{"a failed pass at 11:00 is sent", failedPass,
			[]Sample{sample("LAN", 11, 0, 1)}, []string{"11:00 FIRING failed-pass LAN 1 to morning list"}},
		{"Ystmark isn't held", failedPass,
			[]Sample{sample("YST", 10, 0, 1)}, []string{"10:00 FIRING failed-pass YST 1 to morning list"}},
		{"still firing at 11:00, so paged at 11:00", passFailures,
			[]Sample{sample("LAN", 10, 30, 2), sample("LAN", 10, 45, 2), sample("LAN", 11, 0, 2)},
			[]string{"11:00 FIRING pass-failures LAN 2 to Nkem Halloran"}},
		{"resolved inside the silence, so nothing sent", passFailures,
			[]Sample{sample("LAN", 10, 30, 2), sample("LAN", 10, 50, 0), sample("LAN", 11, 0, 0)}, []string{}},
		{"over by 11:00, so nothing sent", passFailures,
			[]Sample{sample("LAN", 10, 30, 2), sample("LAN", 11, 0, 0)}, []string{}},
		{"a disk full from 10:50 still waits 15 minutes from 10:50", diskFull,
			[]Sample{sample("LAN", 10, 50, 0.92), sample("LAN", 10, 55, 0.92), sample("LAN", 11, 0, 0.92), sample("LAN", 11, 5, 0.92)},
			[]string{"11:05 FIRING disk-full LAN 0.92 to Nkem Halloran"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pager{Rules: []Rule{passFailures, failedPass, diskFull}, Rota: rota, Silences: []Silence{gen}}
			got := []string{}
			for _, s := range tt.samples {
				if a, ok := p.Observe(tt.rule.Name, s); ok {
					got = append(got, a.String()+" to "+a.To)
				}
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("Observe(%q, ...) with the generator's silence sent %q, want %q", tt.rule.Name, got, tt.want)
			}
		})
	}
}

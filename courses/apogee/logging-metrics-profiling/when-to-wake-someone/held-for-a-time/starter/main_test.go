package main

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// diskFull is Ingrid's disk rule from Monday 8 January, and failedPass her rule for each
// failed pass.
var (
	diskFull   = Rule{Name: "disk-full", Severity: Critical, FireAt: 0.90, ClearBelow: 0.85, For: 15 * time.Minute, Repeat: 12 * time.Hour}
	failedPass = Rule{Name: "failed-pass", Severity: Warning, FireAt: 1, ClearBelow: 1}
)

// t0 is when each test's disk first reads 91%: Wednesday 3 January, 09:00.
var t0 = time.Date(2029, 1, 3, 9, 0, 0, 0, time.UTC)

// minutes returns Lantern Atoll's samples at value v, one a minute from minute from to minute
// to after t0, both included.
func minutes(from, to int, v float64) []Sample {
	var ss []Sample
	for m := from; m <= to; m++ {
		ss = append(ss, Sample{At: t0.Add(time.Duration(m) * time.Minute), Station: "LAN", Value: v})
	}
	return ss
}

// observe gives each sample to a new pager holding r, and returns what it sent, as the
// pager's log writes alerts.
func observe(r Rule, samples ...[]Sample) []string {
	p := &Pager{Rules: []Rule{r}}
	sent := []string{}
	for _, s := range slices.Concat(samples...) {
		if a, ok := p.Observe(r.Name, s); ok {
			sent = append(sent, a.String())
		}
	}
	return sent
}

func TestHeldForATime(t *testing.T) {
	tests := []struct {
		name    string
		rule    Rule
		samples [][]Sample
		want    []string
	}{
		{"91% for four minutes while a recording lands", diskFull,
			[][]Sample{minutes(0, 3, 0.91), minutes(4, 4, 0.80)}, []string{}},
		{"91% at 09:00 and at 09:14:59", diskFull,
			[][]Sample{{{At: t0, Station: "LAN", Value: 0.91}, {At: t0.Add(14*time.Minute + 59*time.Second), Station: "LAN", Value: 0.91}}}, []string{}},
		{"91% for exactly 15 minutes", diskFull,
			[][]Sample{minutes(0, 15, 0.91)}, []string{"09:15 FIRING disk-full LAN 0.91"}},
		{"a dip to 89% at 09:10 starts the wait again", diskFull,
			[][]Sample{minutes(0, 9, 0.91), minutes(10, 10, 0.89), minutes(11, 26, 0.91)}, []string{"09:26 FIRING disk-full LAN 0.91"}},
		{"a failed pass, For 0", failedPass,
			[][]Sample{minutes(0, 0, 1)}, []string{"09:00 FIRING failed-pass LAN 1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := observe(tt.rule, tt.samples...); !slices.Equal(got, tt.want) {
				t.Errorf("Observe(%q, ...) sent %s, want %s", tt.rule.Name, list(got), list(tt.want))
			}
		})
	}
}

// list writes alerts for a message: [] when there are none.
func list(alerts []string) string {
	if len(alerts) > 6 {
		return fmt.Sprintf("%q and %d more", alerts[:6], len(alerts)-6)
	}
	return fmt.Sprintf("%q", alerts)
}

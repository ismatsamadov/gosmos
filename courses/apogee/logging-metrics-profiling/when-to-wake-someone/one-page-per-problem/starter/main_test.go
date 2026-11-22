package main

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// passFailures and failedPass are two of Ingrid's rules from Monday 8 January, and
// noReminders is pass-failures without its reminder.
var (
	passFailures = Rule{Name: "pass-failures", Severity: Critical, FireAt: 2, ClearBelow: 1, Repeat: 24 * time.Hour}
	failedPass   = Rule{Name: "failed-pass", Severity: Warning, FireAt: 1, ClearBelow: 1}
	noReminders  = Rule{Name: "pass-failures", Severity: Critical, FireAt: 2, ClearBelow: 1}
)

// failures are the LOS times of Lantern Atoll's six failed passes on the night of 27-28
// December.
var failures = []time.Time{
	time.Date(2028, 12, 27, 19, 41, 28, 0, time.UTC),
	time.Date(2028, 12, 27, 20, 20, 49, 0, time.UTC),
	time.Date(2028, 12, 28, 0, 45, 16, 0, time.UTC),
	time.Date(2028, 12, 28, 1, 57, 49, 0, time.UTC),
	time.Date(2028, 12, 28, 2, 19, 14, 0, time.UTC),
	time.Date(2028, 12, 28, 2, 44, 57, 0, time.UTC),
}

// night returns a rule's samples at Lantern Atoll, every minute from 27 December 19:00 to
// 28 December 09:00: the failed passes with LOS in the window d before, six hours for
// pass-failures and one minute for failed-pass.
func night(d time.Duration) []Sample {
	var ss []Sample
	start := time.Date(2028, 12, 27, 19, 0, 0, 0, time.UTC)
	for at := start; !at.After(start.Add(14 * time.Hour)); at = at.Add(time.Minute) {
		n := 0
		for _, los := range failures {
			if los.After(at.Add(-d)) && !los.After(at) {
				n++
			}
		}
		ss = append(ss, Sample{At: at, Station: "LAN", Value: float64(n)})
	}
	return ss
}

// steady returns samples at value v for each station in turn, every minute from start to
// start+d, both included.
func steady(start time.Time, d time.Duration, v float64, stations ...string) []Sample {
	var ss []Sample
	for at := start; !at.After(start.Add(d)); at = at.Add(time.Minute) {
		for _, st := range stations {
			ss = append(ss, Sample{At: at, Station: st, Value: v})
		}
	}
	return ss
}

// observe gives each sample to a new pager holding r, and returns what it sent, each alert
// as the pager's log writes it after its day.
func observe(r Rule, samples ...[]Sample) []string {
	p := &Pager{Rules: []Rule{r}}
	sent := []string{}
	for _, s := range slices.Concat(samples...) {
		if a, ok := p.Observe(r.Name, s); ok {
			sent = append(sent, a.At.Format("Jan 2 ")+a.String())
		}
	}
	return sent
}

func TestOnePagePerProblem(t *testing.T) {
	jan3 := time.Date(2029, 1, 3, 9, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		rule    Rule
		samples [][]Sample
		want    []string
	}{
		{"the night of 27-28 December", passFailures, [][]Sample{night(6 * time.Hour)},
			[]string{"Dec 27 20:21 FIRING pass-failures LAN 2", "Dec 28 08:45 RESOLVED pass-failures LAN 0"}},
		{"two days failing, reminded after exactly 24 hours, then 24 more", passFailures,
			[][]Sample{steady(jan3, 48*time.Hour, 2, "LAN")},
			[]string{"Jan 3 09:00 FIRING pass-failures LAN 2", "Jan 4 09:00 FIRING pass-failures LAN 2", "Jan 5 09:00 FIRING pass-failures LAN 2"}},
		{"Lantern Atoll and Ystmark failing at once", passFailures,
			[][]Sample{steady(jan3, time.Hour, 2, "LAN", "YST")},
			[]string{"Jan 3 09:00 FIRING pass-failures LAN 2", "Jan 3 09:00 FIRING pass-failures YST 2"}},
		{"Repeat 0 never reminds", noReminders,
			[][]Sample{steady(jan3, 49*time.Hour, 2, "LAN")},
			[]string{"Jan 3 09:00 FIRING pass-failures LAN 2"}},
		{"a Warning that clears sends no RESOLVED", failedPass,
			[][]Sample{night(time.Minute)[:90]},
			[]string{"Dec 27 19:42 FIRING failed-pass LAN 1", "Dec 27 20:21 FIRING failed-pass LAN 1"}},
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

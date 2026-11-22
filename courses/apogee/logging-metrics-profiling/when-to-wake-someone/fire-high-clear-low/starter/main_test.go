package main

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// diskFull is Ingrid's disk rule from Monday 8 January.
var diskFull = Rule{Name: "disk-full", Severity: Critical, FireAt: 0.90, ClearBelow: 0.85, For: 15 * time.Minute, Repeat: 12 * time.Hour}

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

// bobbing returns eight hours of samples from 09:00, a recording landing and being collected
// every half hour: 91% for twenty minutes, then 88% for ten, over and over.
func bobbing() []Sample {
	var ss []Sample
	for m := 0; m < 8*60; m += 30 {
		ss = append(ss, minutes(m, m+19, 0.91)...)
		ss = append(ss, minutes(m+20, m+29, 0.88)...)
	}
	return ss
}

// changes gives each sample to a new pager holding diskFull and returns only the alerts where
// the problem changed: a FIRING that starts it and the RESOLVED that ends it. Until the pager
// sends one page per problem, it repeats FIRING at every sample while a problem lasts.
func changes(samples ...[]Sample) []string {
	p := &Pager{Rules: []Rule{diskFull}}
	got := []string{}
	firing := false
	for _, s := range slices.Concat(samples...) {
		a, ok := p.Observe(diskFull.Name, s)
		if !ok || a.Kind == "FIRING" && firing {
			continue
		}
		firing = a.Kind == "FIRING"
		got = append(got, a.String())
	}
	return got
}

func TestFireHighClearLow(t *testing.T) {
	tests := []struct {
		name    string
		samples [][]Sample
		want    []string
	}{
		{"a disk bobbing between 88% and 91% all night",
			[][]Sample{bobbing()},
			[]string{"09:15 FIRING disk-full LAN 0.91"}},
		{"84% resolves it",
			[][]Sample{minutes(0, 15, 0.91), minutes(16, 20, 0.84)},
			[]string{"09:15 FIRING disk-full LAN 0.91", "09:16 RESOLVED disk-full LAN 0.84"}},
		{"exactly 85% keeps it firing",
			[][]Sample{minutes(0, 15, 0.91), minutes(16, 30, 0.85)},
			[]string{"09:15 FIRING disk-full LAN 0.91"}},
		{"after resolving, a new rise waits out For again",
			[][]Sample{minutes(0, 15, 0.91), minutes(16, 16, 0.84), minutes(17, 32, 0.92)},
			[]string{"09:15 FIRING disk-full LAN 0.91", "09:16 RESOLVED disk-full LAN 0.84", "09:32 FIRING disk-full LAN 0.92"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := changes(tt.samples...); !slices.Equal(got, tt.want) {
				t.Errorf("Observe(%q, ...) fired and resolved %s, want %s", diskFull.Name, list(got), list(tt.want))
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

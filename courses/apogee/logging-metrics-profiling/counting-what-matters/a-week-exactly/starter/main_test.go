package main

import (
	"testing"
	"time"
)

// at returns a time written in RFC 3339, such as 2028-12-13T13:09:43Z.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// readings makes one reading of each value, five minutes apart from Tuesday 12 December 10:00.
func readings(values ...uint64) []Reading {
	var rs []Reading
	for i, v := range values {
		rs = append(rs, Reading{At: at("2028-12-12T10:00:00Z").Add(time.Duration(i) * 5 * time.Minute), Value: v})
	}
	return rs
}

func TestIncrease(t *testing.T) {
	tests := []struct {
		name   string
		values []uint64
		want   uint64
	}{
		{"a restart between 5 and 0", []uint64{3, 5, 0, 1}, 3},
		{"no restart", []uint64{0, 2, 2, 5}, 5},
		{"a restart that has counted 1 already", []uint64{4, 6, 1}, 3},
		{"a counter that never moved", []uint64{7, 7, 7}, 0},
		{"one reading", []uint64{9}, 0},
		{"no readings", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := increase(readings(tt.values...)); got != tt.want {
				t.Errorf("increase(%v) = %d, want %d", tt.values, got, tt.want)
			}
		})
	}
}

// brine3 is Brine-3's 24 passes at Lantern Atoll in Brinewatch's first exact week, by LOS:
// 14926 failed in Wednesday evening's downpour, and every other one flew.
var brine3 = []struct{ id, los, outcome string }{
	{"brine-3/LAN/14911", "2028-12-12T13:39:11Z", "flown"},
	{"brine-3/LAN/14919", "2028-12-13T01:25:10Z", "flown"},
	{"brine-3/LAN/14926", "2028-12-13T13:09:43Z", "failed"},
	{"brine-3/LAN/14934", "2028-12-14T00:55:38Z", "flown"},
	{"brine-3/LAN/14935", "2028-12-14T02:30:36Z", "flown"},
	{"brine-3/LAN/14941", "2028-12-14T12:40:13Z", "flown"},
	{"brine-3/LAN/14949", "2028-12-15T00:26:03Z", "flown"},
	{"brine-3/LAN/14950", "2028-12-15T02:02:09Z", "flown"},
	{"brine-3/LAN/14956", "2028-12-15T12:10:41Z", "flown"},
	{"brine-3/LAN/14957", "2028-12-15T13:45:42Z", "flown"},
	{"brine-3/LAN/14964", "2028-12-15T23:56:25Z", "flown"},
	{"brine-3/LAN/14965", "2028-12-16T01:33:17Z", "flown"},
	{"brine-3/LAN/14971", "2028-12-16T11:41:05Z", "flown"},
	{"brine-3/LAN/14972", "2028-12-16T13:17:13Z", "flown"},
	{"brine-3/LAN/14979", "2028-12-16T23:26:42Z", "flown"},
	{"brine-3/LAN/14980", "2028-12-17T01:04:17Z", "flown"},
	{"brine-3/LAN/14986", "2028-12-17T11:11:26Z", "flown"},
	{"brine-3/LAN/14987", "2028-12-17T12:48:22Z", "flown"},
	{"brine-3/LAN/14994", "2028-12-17T22:56:55Z", "flown"},
	{"brine-3/LAN/14995", "2028-12-18T00:35:10Z", "flown"},
	{"brine-3/LAN/15001", "2028-12-18T10:41:43Z", "flown"},
	{"brine-3/LAN/15002", "2028-12-18T12:19:21Z", "flown"},
	{"brine-3/LAN/15009", "2028-12-18T22:26:59Z", "flown"},
	{"brine-3/LAN/15010", "2028-12-19T00:06:00Z", "flown"},
}

// week returns the pager's readings of passes_total for Brine-3 and outcome at Lantern Atoll,
// every five minutes from Tuesday 12 December 10:00, when the counters arrived at 0, to
// Tuesday 19 December 10:00. A pass counts from its LOS, and the Monday restart at 08:00:30
// on 18 December sets every counter back to 0.
func week(outcome string) []Reading {
	start, restart := at("2028-12-12T10:00:00Z"), at("2028-12-18T08:00:30Z")
	var los []time.Time
	for _, p := range brine3 {
		if p.outcome == outcome {
			los = append(los, at(p.los))
		}
	}
	var rs []Reading
	for t := start; !t.After(start.Add(7 * 24 * time.Hour)); t = t.Add(5 * time.Minute) {
		since := start
		if !t.Before(restart) {
			since = restart
		}
		var n uint64
		for _, l := range los {
			if l.After(since) && !l.After(t) {
				n++
			}
		}
		rs = append(rs, Reading{At: t, Value: n})
	}
	return rs
}

func TestBrinewatchWeek(t *testing.T) {
	for _, c := range []struct {
		outcome string
		want    uint64
	}{
		{"flown", 23},
		{"failed", 1},
	} {
		rs := week(c.outcome)
		if got := increase(rs); got != c.want {
			t.Errorf("increase(the week's %d readings of Brine-3 %s at Lantern Atoll) = %d, want %d", len(rs), c.outcome, got, c.want)
		}
	}
}

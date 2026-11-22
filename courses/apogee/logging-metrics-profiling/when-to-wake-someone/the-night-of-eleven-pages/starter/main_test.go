package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
)

// TestReplay plays the night of 27-28 December through the new rules, as Ingrid does.
func TestReplay(t *testing.T) {
	got := []string{}
	for _, a := range replay(newRules()) {
		got = append(got, a.String()+" to "+a.To)
	}
	want := []string{
		"19:42 FIRING failed-pass LAN 1 to morning list",
		"20:21 FIRING pass-failures LAN 2 to Tomás Ferreira-Lund",
		"20:21 FIRING failed-pass LAN 1 to morning list",
		"00:46 FIRING failed-pass LAN 1 to morning list",
		"01:58 FIRING failed-pass LAN 1 to morning list",
		"02:20 FIRING failed-pass LAN 1 to morning list",
		"02:45 FIRING failed-pass LAN 1 to morning list",
		"08:45 RESOLVED pass-failures LAN 0 to Tomás Ferreira-Lund",
	}
	if !slices.Equal(got, want) {
		if len(got) > 20 {
			got = append(got[:20], fmt.Sprintf("and %d more", len(got)-20))
		}
		t.Errorf("replay(newRules()) sent:\n\t%s\nwant:\n\t%s", strings.Join(got, "\n\t"), strings.Join(want, "\n\t"))
	}
}

// describe writes a rule's settings for a message.
func describe(r Rule) string {
	sev := map[Severity]string{Info: "Info", Warning: "Warning", Critical: "Critical"}[r.Severity]
	return fmt.Sprintf("FireAt %g, ClearBelow %g, For %v, Repeat %v, %s", r.FireAt, r.ClearBelow, r.For, r.Repeat, sev)
}

func TestNewRules(t *testing.T) {
	want := []Rule{
		{Name: "pass-failures", Severity: Critical, FireAt: 2, ClearBelow: 1, Repeat: 24 * time.Hour},
		{Name: "failed-pass", Severity: Warning, FireAt: 1, ClearBelow: 1},
		{Name: "disk-full", Severity: Critical, FireAt: 0.90, ClearBelow: 0.85, For: 15 * time.Minute, Repeat: 12 * time.Hour},
	}
	got := newRules()
	for _, w := range want {
		i := slices.IndexFunc(got, func(r Rule) bool { return r.Name == w.Name })
		switch {
		case i < 0:
			t.Errorf("newRules() has no %q, want %s", w.Name, describe(w))
		case got[i] != w:
			t.Errorf("newRules()'s %q = %s, want %s", w.Name, describe(got[i]), describe(w))
		}
	}
	if len(got) != len(want) {
		t.Errorf("newRules() has %d rules, want 3: pass-failures, failed-pass and disk-full", len(got))
	}
}

// oldPages is what the old rules sent for the night, ten seconds after each LOS:
// lan-failed-pass for every failed pass, and pass-failures while at least 3 of the station's
// last 15 passes had failed, until Tomás muted Lantern Atoll at 07:28.
func oldPages(results []Result) []string {
	mute := time.Date(2028, 12, 28, 7, 28, 0, 0, time.UTC)
	var pages []string
	for i, r := range results {
		at := r.LOS.Add(10 * time.Second)
		if !at.Before(mute) {
			break
		}
		if r.Outcome == "failed" {
			pages = append(pages, at.Format("15:04:05")+" lan-failed-pass")
		}
		failed := 0
		for _, q := range results[max(0, i-14) : i+1] {
			if q.Outcome == "failed" {
				failed++
			}
		}
		if failed >= 3 {
			pages = append(pages, at.Format("15:04:05")+" pass-failures")
		}
	}
	return pages
}

// TestOldRules checks the replay's night against Tomás's count: the old rules, on the same
// results, paged eleven times.
func TestOldRules(t *testing.T) {
	if got := oldPages(night); len(got) != 11 {
		t.Errorf("the old rules on the night = %d pages %q, want 11", len(got), got)
	}
}

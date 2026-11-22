package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"testing"
)

// passesFlags makes the passes command's flag set, as gctl defines it.
func passesFlags() (*flag.FlagSet, *string) {
	fs := newFlagSet("passes", "passes [-station CODE] file")
	return fs, fs.String("station", "", "show only the bookings at station `CODE`")
}

// holds reports whether text holds every one of wants, or is empty when there are none.
func holds(text string, wants []string) bool {
	if len(wants) == 0 {
		return text == ""
	}
	for _, w := range wants {
		if !strings.Contains(text, w) {
			return false
		}
	}
	return true
}

// describe says what holds checks for.
func describe(wants []string) string {
	if len(wants) == 0 {
		return "nothing"
	}
	return fmt.Sprintf("text holding %q", wants)
}

func TestParseFlags(t *testing.T) {
	for _, tc := range []struct {
		name           string
		args           []string
		status         int
		ok             bool
		station        string
		stdout, stderr []string
	}{
		{"Salar Quieto", []string{"-station", "SQT", "week.json"}, 0, true, "SQT", nil, nil},
		{"no flags", []string{"week.json"}, 0, true, "", nil, nil},
		{"short help", []string{"-h"}, 0, false, "", []string{"usage: gctl passes"}, nil},
		{"long help", []string{"--help"}, 0, false, "", []string{"usage: gctl passes"}, nil},
		{"misspelled flag", []string{"-stn", "SQT"}, 2, false, "", nil, []string{"-stn", "usage: gctl passes"}},
		{"no value", []string{"-station"}, 2, false, "", nil, []string{"-station", "usage: gctl passes"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs, station := passesFlags()
			var stdout, stderr bytes.Buffer
			status, ok := parseFlags(fs, tc.args, &stdout, &stderr)
			if status != tc.status || ok != tc.ok {
				t.Errorf("parseFlags(passes, %q) = %d, %v, want %d, %v", tc.args, status, ok, tc.status, tc.ok)
			}
			if !holds(stdout.String(), tc.stdout) {
				t.Errorf("parseFlags(passes, %q) stdout = %q, want %s", tc.args, stdout.String(), describe(tc.stdout))
			}
			if !holds(stderr.String(), tc.stderr) {
				t.Errorf("parseFlags(passes, %q) stderr = %q, want %s", tc.args, stderr.String(), describe(tc.stderr))
			}
			if *station != tc.station {
				t.Errorf("parseFlags(passes, %q) set -station to %q, want %q", tc.args, *station, tc.station)
			}
		})
	}
}

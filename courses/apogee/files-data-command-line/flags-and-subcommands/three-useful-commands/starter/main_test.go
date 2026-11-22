package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// week is the recorder's booking log for Tuesday 5 October, out of order.
const week = `{"id":"bk-3203","sat_id":"furrow-a","station":"SQT","aos":"2027-10-05T14:02:10Z"}
{"id":"bk-3201","sat_id":"pipit-2","station":"PEL","aos":"2027-10-05T09:40:00Z"}
{"id":"bk-3202","sat_id":"cirrine-1","station":"SQT","aos":"2027-10-05T06:15:30Z"}
`

// matches reports whether got is empty for an empty want, exactly a want ending in a newline, or holds want.
func matches(got, want string) bool {
	switch {
	case want == "":
		return got == ""
	case strings.HasSuffix(want, "\n"):
		return got == want
	}
	return strings.Contains(got, want)
}

// describe says what matches checks for.
func describe(want string) string {
	switch {
	case want == "":
		return "nothing"
	case strings.HasSuffix(want, "\n"):
		return fmt.Sprintf("%q", want)
	}
	return fmt.Sprintf("text holding %q", want)
}

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		name           string
		args           []string
		status         int
		stdout, stderr string
	}{
		{"stow for wind", []string{"stow", "-reason", "wind 75 km/h", "CWR"}, 0, "CWR stowed: wind 75 km/h\n", ""},
		{"stow without a reason", []string{"stow", "SQT"}, 0, "SQT stowed: no reason given\n", ""},
		{"help for stow", []string{"stow", "-h"}, 0, "usage: gctl stow [-reason TEXT] CODE", ""},
		{"stow without a code", []string{"stow", "-reason", "dust"}, 2, "", "usage: gctl stow [-reason TEXT] CODE"},
		{"stow with two codes", []string{"stow", "SQT", "CWR"}, 2, "", "usage: gctl stow [-reason TEXT] CODE"},
		{"reason after the code", []string{"stow", "SQT", "-reason", "dust"}, 2, "", "usage: gctl stow [-reason TEXT] CODE"},
		{"misspelled flag", []string{"stow", "-raeson", "wind", "CWR"}, 2, "", "-raeson"},
		{"unknown station", []string{"stow", "XYZ"}, 1, "", "gctl: unknown station \"XYZ\"\n"},
		{"passes at SQT", []string{"passes", "-station", "SQT"}, 0, "2027-10-05T06:15:30Z SQT cirrine-1 bk-3202\n2027-10-05T14:02:10Z SQT furrow-a bk-3203\n", ""},
		{"station table", []string{"stations"}, 0, "YST Ystmark\nPEL Port Ellery\nSQT Salar Quieto\nLAN Lantern Atoll\nCWR Cape Wren\n", ""},
		{"unknown command", []string{"park", "SQT"}, 2, "", "unknown command \"park\""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			status := run(tc.args, strings.NewReader(week), &stdout, &stderr)
			if status != tc.status {
				t.Errorf("run(%q) = %d, want %d", tc.args, status, tc.status)
			}
			if !matches(stdout.String(), tc.stdout) {
				t.Errorf("run(%q) stdout = %q, want %s", tc.args, stdout.String(), describe(tc.stdout))
			}
			if !matches(stderr.String(), tc.stderr) {
				t.Errorf("run(%q) stderr = %q, want %s", tc.args, stderr.String(), describe(tc.stderr))
			}
		})
	}
}

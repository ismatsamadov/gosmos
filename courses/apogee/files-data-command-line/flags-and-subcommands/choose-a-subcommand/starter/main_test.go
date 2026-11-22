package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// week is a booking log of three passes, out of order.
const week = `{"id":"bk-3203","sat_id":"furrow-a","station":"SQT","aos":"2027-10-05T14:02:10Z","los":"2027-10-05T14:11:52Z","note":"priority"}
{"id":"bk-3201","sat_id":"pipit-2","station":"PEL","aos":"2027-10-05T09:40:00Z","los":"2027-10-05T09:49:00Z"}
{"id":"bk-3202","sat_id":"cirrine-1","station":"SQT","aos":"2027-10-05T06:15:30Z","los":"2027-10-05T06:26:00Z"}
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
	path := filepath.Join(t.TempDir(), "week.json")
	if err := os.WriteFile(path, []byte(week), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name           string
		args           []string
		status         int
		stdout, stderr string
	}{
		{"station table", []string{"stations"}, 0, "YST Ystmark\nPEL Port Ellery\nSQT Salar Quieto\nLAN Lantern Atoll\nCWR Cape Wren\n", ""},
		{"passes at SQT", []string{"passes", "-station", "SQT", "week.json"}, 0, "2027-10-05T06:15:30Z SQT cirrine-1 bk-3202\n2027-10-05T14:02:10Z SQT furrow-a bk-3203\n", ""},
		{"help for passes", []string{"passes", "-h"}, 0, "usage: gctl passes", ""},
		{"gctl help", []string{"-h"}, 0, "usage: gctl <command>", ""},
		{"no command", nil, 2, "", "gctl needs a command\nusage: gctl <command>"},
		{"unknown command", []string{"fly"}, 2, "", "unknown command \"fly\"\nusage: gctl <command>"},
		{"flag before the command", []string{"-station", "SQT", "passes"}, 2, "", "unknown command \"-station\"\nusage: gctl <command>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := slices.Clone(tc.args)
			if i := slices.Index(args, "week.json"); i >= 0 {
				args[i] = path
			}
			var stdout, stderr bytes.Buffer
			status := run(args, strings.NewReader(""), &stdout, &stderr)
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

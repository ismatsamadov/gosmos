package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReport(t *testing.T) {
	for _, tc := range []struct {
		name           string
		lines          []string
		err            error
		status         int
		stdout, stderr string
	}{
		{"two passes", []string{"2027-10-05T06:15:30Z SQT cirrine-1 bk-3202", "2027-10-05T14:02:10Z SQT furrow-a bk-3203"}, nil, 0, "2027-10-05T06:15:30Z SQT cirrine-1 bk-3202\n2027-10-05T14:02:10Z SQT furrow-a bk-3203\n", ""},
		{"empty report", nil, nil, 0, "", ""},
		{"failure", nil, errors.New("open week.json: no such file or directory"), 1, "", "gctl: open week.json: no such file or directory\n"},
		{"failure after some passes", []string{"2027-10-05T06:15:30Z SQT cirrine-1 bk-3202"}, errors.New("booking 2: unexpected EOF"), 1, "", "gctl: booking 2: unexpected EOF\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			status := report(tc.lines, tc.err, &stdout, &stderr)
			if status != tc.status || stdout.String() != tc.stdout || stderr.String() != tc.stderr {
				t.Errorf("report(%q, %v) = %d with stdout %q and stderr %q, want %d with stdout %q and stderr %q", tc.lines, tc.err, status, stdout.String(), stderr.String(), tc.status, tc.stdout, tc.stderr)
			}
		})
	}
}

// twoPasses is a booking log of two passes at Salar Quieto, out of order.
const twoPasses = `{"id":"bk-3203","sat_id":"furrow-a","station":"SQT","aos":"2027-10-05T14:02:10Z"}
{"id":"bk-3202","sat_id":"cirrine-1","station":"SQT","aos":"2027-10-05T06:15:30Z"}
`

// broken is the same log with a third booking whose AOS has no offset.
const broken = twoPasses + `{"id":"bk-3204","sat_id":"pipit-1","station":"SQT","aos":"2027-10-05T21:30:00"}
`

func TestRunPasses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "week.json")
	if err := os.WriteFile(path, []byte(broken), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		args   []string
		stdin  string
		status int
		stdout string
		stderr string // the start of the one line stderr must hold, or empty for nothing
	}{
		{"two passes from the pipe", []string{"-station", "SQT"}, twoPasses, 0, "2027-10-05T06:15:30Z SQT cirrine-1 bk-3202\n2027-10-05T14:02:10Z SQT furrow-a bk-3203\n", ""},
		{"empty pipe", nil, "", 0, "", ""},
		{"no bookings at LAN", []string{"-station", "LAN"}, twoPasses, 0, "", ""},
		{"bad third booking", []string{"week.json"}, "", 1, "", "gctl: booking 3: "},
		{"missing file", []string{"missing.json"}, "", 1, "", "gctl: open "},
		{"unknown station", []string{"-station", "SQX", "week.json"}, "", 1, "", "gctl: unknown station \"SQX\""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var args []string
			for _, a := range tc.args {
				if strings.HasSuffix(a, ".json") {
					a = filepath.Join(dir, a)
				}
				args = append(args, a)
			}
			var stdout, stderr bytes.Buffer
			status := runPasses(args, strings.NewReader(tc.stdin), &stdout, &stderr)
			if status != tc.status {
				t.Errorf("runPasses(%q) = %d, want %d", tc.args, status, tc.status)
			}
			if stdout.String() != tc.stdout {
				t.Errorf("runPasses(%q) stdout = %q, want %q", tc.args, stdout.String(), tc.stdout)
			}
			got := stderr.String()
			if tc.stderr == "" && got != "" {
				t.Errorf("runPasses(%q) stderr = %q, want nothing", tc.args, got)
			}
			if tc.stderr != "" && (!strings.HasPrefix(got, tc.stderr) || strings.Count(got, "\n") != 1 || !strings.HasSuffix(got, "\n")) {
				t.Errorf("runPasses(%q) stderr = %q, want one line starting %q", tc.args, got, tc.stderr)
			}
		})
	}
}

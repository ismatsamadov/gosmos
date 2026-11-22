package main

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// settings is Thursday's settings file.
const settings = `{"station":"SQT","day":"2027-10-21"}`

// export is Yusra's export for Thursday: two passes at SQT, one at SQT on Friday.
const export = `[
{"id":"bk-3410","sat_id":"furrow-a","station":"SQT","aos":"2027-10-21T14:02:10Z","los":"2027-10-21T14:13:52Z"},
{"id":"bk-3413","sat_id":"brine-3","station":"SQT","aos":"2027-10-21T15:40:05Z","los":"2027-10-21T15:51:31Z"},
{"id":"bk-3422","sat_id":"furrow-b","station":"SQT","aos":"2027-10-22T16:05:00Z","los":"2027-10-22T16:16:41Z"}
]`

// typo is the same export with Friday's pass booked at a station that does not exist.
const typo = `[
{"id":"bk-3410","sat_id":"furrow-a","station":"SQT","aos":"2027-10-21T14:02:10Z","los":"2027-10-21T14:13:52Z"},
{"id":"bk-3413","sat_id":"brine-3","station":"SQT","aos":"2027-10-21T15:40:05Z","los":"2027-10-21T15:51:31Z"},
{"id":"bk-3422","sat_id":"furrow-b","station":"SQX","aos":"2027-10-22T16:05:00Z","los":"2027-10-22T16:16:41Z"}
]`

// stick writes files into a new folder, as they would be on the memory stick, and returns the folder.
func stick(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestPlanDay(t *testing.T) {
	plan, err := planDay(stick(t, map[string]string{"settings.json": settings, "bookings.json": export}), "json")
	if plan.Day != "2027-10-21" || len(plan.Passes) != 2 || err != nil {
		t.Errorf("planDay(good stick) = %+v, %v, want 2 passes on 2027-10-21, nil", plan, err)
	}
}

func TestPlanDayErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		files  map[string]string
		format string
		prefix string // how the message must start
		target error  // nil when any error will do
	}{
		{"no settings on the stick", map[string]string{"bookings.json": export}, "json", "settings: ", fs.ErrNotExist},
		{"misspelled setting", map[string]string{"settings.json": `{"day":"2027-10-21","min_elev":5}`, "bookings.json": export}, "json", "settings: ", nil},
		{"no booking sheet", map[string]string{"settings.json": settings, "bookings.json": export}, "csv", "bookings: ", fs.ErrNotExist},
		{"export cut off halfway", map[string]string{"settings.json": settings, "bookings.json": export[:120]}, "json", "bookings: ", io.ErrUnexpectedEOF},
		{"empty export", map[string]string{"settings.json": settings, "bookings.json": ""}, "json", "bookings: ", io.EOF},
		{"station SQX", map[string]string{"settings.json": settings, "bookings.json": typo}, "json", "bookings: ", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := planDay(stick(t, tc.files), tc.format)
			switch {
			case err == nil:
				t.Errorf("planDay(%s) = a plan of %d passes, nil, want an error", tc.name, len(plan.Passes))
			case !strings.HasPrefix(err.Error(), tc.prefix):
				t.Errorf("planDay(%s) error = %q, want it to start with %q", tc.name, err, tc.prefix)
			case tc.target != nil && !errors.Is(err, tc.target):
				t.Errorf("planDay(%s) error = %v, want it to wrap %v", tc.name, err, tc.target)
			}
		})
	}
}

func TestFailuresEndToEnd(t *testing.T) {
	good := map[string]string{"settings.json": settings, "bookings.json": export}
	full := errors.New("no space left on the stick")
	for _, tc := range []struct {
		name   string
		args   []string // the folder is added last
		files  map[string]string
		output io.Writer // nil for a working standard output
		code   int
		stderr string // how standard error must start
	}{
		{"bad flag", []string{"day", "-speed", "3"}, good, nil, 2, "flag provided but not defined: -speed\n"},
		{"bad format", []string{"day", "-format", "xml"}, good, nil, 2, "-format \"xml\": want json or csv\n"},
		{"bad data", []string{"day"}, map[string]string{"settings.json": settings, "bookings.json": typo}, nil, 1, "gctl: bookings: "},
		{"missing settings", []string{"day"}, map[string]string{"bookings.json": export}, nil, 1, "gctl: settings: "},
		{"broken stick", []string{"day"}, map[string]string{"settings.json": settings, "bookings.json": export[:120]}, nil, 1, "gctl: bookings: "},
		{"failed output", []string{"day"}, good, writerFunc(func([]byte) (int, error) { return 0, full }), 1, "gctl: " + full.Error()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			out := tc.output
			if out == nil {
				out = &stdout
			}
			code := run(append(tc.args, stick(t, tc.files)), strings.NewReader(""), out, &stderr)
			if code != tc.code {
				t.Errorf("run(%s) = %d, want %d", tc.name, code, tc.code)
			}
			if stdout.Len() > 0 {
				t.Errorf("run(%s) printed %q, want nothing on standard output", tc.name, stdout.String())
			}
			if !strings.HasPrefix(stderr.String(), tc.stderr) {
				t.Errorf("run(%s) standard error = %q, want it to start with %q", tc.name, stderr.String(), tc.stderr)
			}
			if tc.code == 1 && strings.Count(stderr.String(), "\n") != 1 {
				t.Errorf("run(%s) standard error = %q, want exactly one line", tc.name, stderr.String())
			}
			if tc.code == 2 && !strings.HasSuffix(stderr.String(), dayUsage) {
				t.Errorf("run(%s) standard error = %q, want it to end with the day command's usage", tc.name, stderr.String())
			}
		})
	}
}

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the test function.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

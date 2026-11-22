package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// copyStick copies testdata/stick into a new folder, the stick as Paz plugs it in, and returns the folder.
func copyStick(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS(filepath.Join("testdata", "stick"))); err != nil {
		t.Fatal(err)
	}
	return dir
}

// emptyHome moves the test into an empty folder and returns it: gctl must leave it empty.
func emptyHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Chdir(home)
	return home
}

// names returns the names a folder holds.
func names(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range entries {
		out = append(out, e.Name())
	}
	return out
}

// firstDifference returns the number of the first line where got and want differ, and both lines.
func firstDifference(got, want string) (int, string, string) {
	g, w := strings.SplitAfter(got, "\n"), strings.SplitAfter(want, "\n")
	for i := range max(len(g), len(w)) {
		gotLine, wantLine := "(nothing)", "(nothing)"
		if i < len(g) {
			gotLine = g[i]
		}
		if i < len(w) {
			wantLine = w[i]
		}
		if gotLine != wantLine {
			return i + 1, gotLine, wantLine
		}
	}
	return 0, "", ""
}

func TestCalibrationDay(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "thursday.golden"))
	if err != nil {
		t.Fatal(err)
	}
	golden := string(data)
	for _, format := range []string{"json", "csv"} {
		t.Run(format, func(t *testing.T) {
			dir := copyStick(t)
			home := emptyHome(t)
			for range 2 { // the second run replaces the first report
				var stdout, stderr bytes.Buffer
				code := run([]string{"day", "-format", format, dir}, strings.NewReader(""), &stdout, &stderr)
				if code != 0 || stderr.Len() > 0 {
					t.Fatalf("run(day -format %s) = %d with standard error %q, want 0 and nothing", format, code, stderr.String())
				}
				if n, got, want := firstDifference(stdout.String(), golden); n > 0 {
					t.Fatalf("run(day -format %s) printed line %d = %q, want %q", format, n, got, want)
				}
			}
			if got := names(t, home); len(got) > 0 {
				t.Fatalf("the current folder gained %q, want every file inside the day's folder", got)
			}
			daily := filepath.Join(dir, "reports", "daily")
			saved, err := os.ReadFile(filepath.Join(daily, "2027-10-21.txt"))
			if err != nil {
				t.Fatalf("reports/daily/2027-10-21.txt in the day's folder: %v", err)
			}
			if n, got, want := firstDifference(string(saved), golden); n > 0 {
				t.Errorf("saved report line %d = %q, want %q", n, got, want)
			}
			if got := names(t, daily); !slices.Equal(got, []string{"2027-10-21.txt"}) {
				t.Errorf("reports/daily holds %q, want only 2027-10-21.txt", got)
			}
		})
	}
}

func TestBadBookingKeepsLastNightsPlan(t *testing.T) {
	dir := copyStick(t)
	sheet, err := os.ReadFile(filepath.Join(dir, "bookings.csv"))
	if err != nil {
		t.Fatal(err)
	}
	typo := strings.Replace(string(sheet), "SQT,bk-3422,", "SQX,bk-3422,", 1)
	daily := filepath.Join(dir, "reports", "daily")
	last := "Calibration day 2027-10-21 at SQT (Salar Quieto)\nBooked passes (UTC): 3\n"
	if err := os.WriteFile(filepath.Join(dir, "bookings.csv"), []byte(typo), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(daily, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(daily, "2027-10-21.txt"), []byte(last), 0600); err != nil {
		t.Fatal(err)
	}
	home := emptyHome(t)
	var stdout, stderr bytes.Buffer
	code := run([]string{"day", "-format", "csv", dir}, strings.NewReader(""), &stdout, &stderr)
	if code != 1 || stdout.Len() > 0 || !strings.HasPrefix(stderr.String(), "gctl: ") {
		t.Errorf("run(sheet with SQX) = %d, printed %q, standard error %q, want 1, nothing, a gctl: line", code, stdout.String(), stderr.String())
	}
	saved, err := os.ReadFile(filepath.Join(daily, "2027-10-21.txt"))
	if string(saved) != last || err != nil {
		t.Errorf("plan after a failed day = %q, %v, want last night's %q, nil", saved, err, last)
	}
	if got := names(t, home); len(got) > 0 {
		t.Errorf("the current folder gained %q, want it left alone", got)
	}
}

func TestNothingPrintedWhenTheSaveFails(t *testing.T) {
	dir := copyStick(t)
	if err := os.WriteFile(filepath.Join(dir, "reports"), []byte("a file where the folder should be\n"), 0600); err != nil {
		t.Fatal(err)
	}
	emptyHome(t)
	var stdout, stderr bytes.Buffer
	code := run([]string{"day", dir}, strings.NewReader(""), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stderr.String(), "gctl: ") {
		t.Errorf("run(reports blocked) = %d, standard error %q, want 1 and a gctl: line", code, stderr.String())
	}
	if stdout.Len() > 0 {
		t.Errorf("run(reports blocked) printed %q, want nothing: a plan that was not saved is not printed", stdout.String())
	}
}

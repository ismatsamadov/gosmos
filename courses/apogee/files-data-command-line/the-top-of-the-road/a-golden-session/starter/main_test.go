package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// transcript writes one session as the golden files hold it: standard output, standard error, exit status.
func transcript(stdout, stderr string, code int) string {
	return fmt.Sprintf("[stdout]\n%s[stderr]\n%s[exit %d]\n", stdout, stderr, code)
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

func TestGoldenSessions(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"export", []string{"day", "testdata/stick"}},
		{"sheet", []string{"day", "-format", "csv", "testdata/stick"}},
		{"no-folder", []string{"day"}},
		{"day-help", []string{"day", "-h"}},
		{"help", []string{"-h"}},
		{"unknown", []string{"park", "SQT"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want, err := os.ReadFile(filepath.Join("testdata", tc.name+".golden"))
			if err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code := run(tc.args, strings.NewReader(""), &stdout, &stderr)
			got := transcript(stdout.String(), stderr.String(), code)
			if n, gotLine, wantLine := firstDifference(got, string(want)); n > 0 {
				t.Errorf("run(%q) line %d = %q, want %q, as testdata/%s.golden says", tc.args, n, gotLine, wantLine, tc.name)
			}
		})
	}
}

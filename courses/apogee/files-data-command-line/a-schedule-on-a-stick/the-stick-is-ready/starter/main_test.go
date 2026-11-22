package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportDay(t *testing.T) {
	root := filepath.Join(t.TempDir(), "van")
	for _, tc := range []struct{ day, text string }{{"2027-07-12", "PEL first\n"}, {"2027-07-13", "SQT second\n"}, {"2027-07-12", "replacement\n"}} {
		if err := exportDay(root, tc.day, tc.text); err != nil {
			t.Fatalf("exportDay(%q) error = %v, want nil", tc.day, err)
		}
		got, err := os.ReadFile(filepath.Join(root, "reports", "daily", tc.day+".txt"))
		if err != nil || string(got) != tc.text {
			t.Errorf("exportDay(%q) saved = %q, %v, want %q, nil", tc.day, got, err, tc.text)
		}
	}
	blocked := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocked, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := exportDay(blocked, "day", "new"); err == nil {
		t.Error("exportDay(blocked root) error = nil, want an error")
	}
}

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveReport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.txt")
	for _, text := range []string{"Long report\n", "PEL\n", ""} {
		if err := saveReport(path, text); err != nil {
			t.Fatalf("saveReport(path, text) error = %v, want nil", err)
		}
		got, err := os.ReadFile(path)
		if err != nil || string(got) != text {
			t.Errorf("saved bytes = %q, %v, want %q, nil", got, err, text)
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0077 != 0 {
		t.Errorf("new file mode = %o, want no group or other bits", info.Mode().Perm())
	}
	if err := saveReport(filepath.Join(t.TempDir(), "missing", "file"), "text"); err == nil {
		t.Error("saveReport(missing parent) error = nil, want an error")
	}
}

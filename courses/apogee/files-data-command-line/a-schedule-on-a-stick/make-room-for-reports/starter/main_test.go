package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareReports(t *testing.T) {
	root := filepath.Join(t.TempDir(), "van")
	for range 2 {
		if err := prepareReports(root); err != nil {
			t.Fatalf("prepareReports(root) error = %v, want nil", err)
		}
		info, err := os.Stat(filepath.Join(root, "reports", "daily"))
		if err != nil || !info.IsDir() {
			t.Fatalf("reports/daily directory = %v, %v, want directory, nil", info, err)
		}
		if info.Mode().Perm()&0077 != 0 {
			t.Errorf("reports/daily mode = %o, want no group or other bits", info.Mode().Perm())
		}
	}
	blocked := t.TempDir()
	if err := os.WriteFile(filepath.Join(blocked, "reports"), []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := prepareReports(blocked); err == nil {
		t.Error("prepareReports(blocked) error = nil, want an error")
	}
}

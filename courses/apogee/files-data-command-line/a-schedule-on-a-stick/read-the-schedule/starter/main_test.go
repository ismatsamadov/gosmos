package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestReadSchedule(t *testing.T) {
	for _, tc := range []struct{ name, text string }{{"passes", "pipit-1 PEL 04:12\n"}, {"empty", ""}, {"spaces", "  PEL\r\n"}} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "passes.txt")
			if err := os.WriteFile(path, []byte(tc.text), 0600); err != nil {
				t.Fatal(err)
			}
			got, err := readSchedule(path)
			if got != tc.text || err != nil {
				t.Errorf("readSchedule(path) = %q, %v, want %q, nil", got, err, tc.text)
			}
		})
	}
	path := filepath.Join(t.TempDir(), "missing.txt")
	got, err := readSchedule(path)
	if got != "" || !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("readSchedule(missing) = %q, %v, want empty and not-exist error", got, err)
	}
}

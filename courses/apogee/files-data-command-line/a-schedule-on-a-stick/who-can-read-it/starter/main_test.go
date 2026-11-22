package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestOwnerOnly(t *testing.T) {
	for _, tc := range []struct {
		name string
		mode fs.FileMode
		want bool
	}{{"owner", 0600, true}, {"group", 0640, false}, {"others", 0644, false}, {"execute", 0601, false}} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "report")
			if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, tc.mode); err != nil {
				t.Fatal(err)
			}
			got, err := ownerOnly(path)
			if got != tc.want || err != nil {
				t.Errorf("ownerOnly(mode %o) = %v, %v, want %v, nil", tc.mode, got, err, tc.want)
			}
		})
	}
	got, err := ownerOnly(filepath.Join(t.TempDir(), "missing"))
	if got || !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("ownerOnly(missing) = %v, %v, want false and not-exist error", got, err)
	}
}

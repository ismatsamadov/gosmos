package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// loadFixture returns the non-empty lines of the file testdata/<name>.
func loadFixture(t testing.TB, name string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

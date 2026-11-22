package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceSchedule(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "passes.txt")
	for _, text := range []string{"old complete schedule\n", "new\n", ""} {
		if err := replaceSchedule(path, []byte(text)); err != nil {
			t.Fatalf("replaceSchedule(path) error = %v, want nil", err)
		}
		got, err := os.ReadFile(path)
		if err != nil || string(got) != text {
			t.Errorf("saved schedule = %q, %v, want %q, nil", got, err, text)
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 {
			t.Errorf("directory entries = %d, want 1 without temporary files", len(entries))
		}
	}
}

func TestReplaceScheduleKeepsOldFileForReaders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "passes.txt")
	old := "old complete schedule\n"
	if err := os.WriteFile(path, []byte(old), 0600); err != nil {
		t.Fatal(err)
	}
	reader, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if err := replaceSchedule(path, []byte("new\n")); err != nil {
		t.Fatalf("replaceSchedule(path) error = %v, want nil", err)
	}
	got, err := io.ReadAll(reader)
	if err != nil || string(got) != old {
		t.Errorf("reader opened before replaceSchedule(path) = %q, %v, want old %q, nil", got, err, old)
	}
}

func TestReplaceScheduleFailureCleansUp(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "blocked")
	if err := os.Mkdir(target, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "keep"), []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := replaceSchedule(target, []byte("new")); err == nil {
		t.Error("replaceSchedule(directory) error = nil, want an error")
	}
	got, err := os.ReadFile(filepath.Join(target, "keep"))
	if err != nil || string(got) != "old" {
		t.Errorf("original contents = %q, %v, want old, nil", got, err)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("entries after failure = %d, want 1 without temporary files", len(entries))
	}
	if err := replaceSchedule(filepath.Join(root, "missing", "file"), nil); err == nil {
		t.Error("replaceSchedule(missing parent) error = nil, want an error")
	}
}

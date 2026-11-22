package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// text returns a report function that writes s.
func text(s string) func(io.Writer) error {
	return func(w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	}
}

// entries returns how many names a folder holds.
func entries(t *testing.T, dir string) int {
	t.Helper()
	list, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	return len(list)
}

func TestSaveDay(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "2027-10-21.txt")
	for _, report := range []string{"Last night's plan: four passes\n", "This morning's plan\n", ""} {
		if err := saveDay(path, text(report)); err != nil {
			t.Fatalf("saveDay(path, %q) = %v, want nil", report, err)
		}
		got, err := os.ReadFile(path)
		if string(got) != report || err != nil {
			t.Errorf("saved report = %q, %v, want %q, nil", got, err, report)
		}
		if n := entries(t, dir); n != 1 {
			t.Errorf("folder holds %d names after saveDay, want 1: no temporary file left behind", n)
		}
	}
}

func TestSaveDayKeepsLastNightsPlan(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "2027-10-21.txt")
	last := "Last night's plan: four passes\n"
	if err := os.WriteFile(path, []byte(last), 0600); err != nil {
		t.Fatal(err)
	}
	problem := errors.New("report stopped at line 3")
	err := saveDay(path, func(w io.Writer) error {
		if _, err := io.WriteString(w, "This morning's plan\nline 2\n"); err != nil {
			return err
		}
		return problem
	})
	if !errors.Is(err, problem) {
		t.Errorf("saveDay(report that fails at line 3) = %v, want the report's error", err)
	}
	got, err := os.ReadFile(path)
	if string(got) != last || err != nil {
		t.Errorf("plan after a failed save = %q, %v, want last night's %q, nil", got, err, last)
	}
	if n := entries(t, dir); n != 1 {
		t.Errorf("folder holds %d names after a failed save, want 1: no temporary file left behind", n)
	}
}

func TestSaveDayLeavesOpenReadersTheOldPlan(t *testing.T) {
	path := filepath.Join(t.TempDir(), "2027-10-21.txt")
	last := "Last night's plan: four passes\n"
	if err := os.WriteFile(path, []byte(last), 0600); err != nil {
		t.Fatal(err)
	}
	reader, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if err := saveDay(path, text("This morning's plan\n")); err != nil {
		t.Fatalf("saveDay(path, report) = %v, want nil", err)
	}
	got, err := io.ReadAll(reader)
	if string(got) != last || err != nil {
		t.Errorf("reader opened before saveDay(path) = %q, %v, want last night's %q, nil: the new report replaces the name, not the old file", got, err, last)
	}
}

func TestSaveDayRenameFails(t *testing.T) {
	dir := t.TempDir()
	blocked := filepath.Join(dir, "2027-10-21.txt")
	if err := os.Mkdir(blocked, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(blocked, "keep"), []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := saveDay(blocked, text("This morning's plan\n")); err == nil {
		t.Errorf("saveDay(a folder's path) = nil, want an error")
	}
	if n := entries(t, dir); n != 1 {
		t.Errorf("folder holds %d names after a failed rename, want 1: no temporary file left behind", n)
	}
	if err := saveDay(filepath.Join(dir, "missing", "2027-10-21.txt"), text("plan\n")); err == nil {
		t.Errorf("saveDay(path in a missing folder) = nil, want an error")
	}
}

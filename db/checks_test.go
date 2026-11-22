package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The cases in testdata/checks are shared with the app's ValidationTests: each is a library
// broken for one check, and the problems the checks report for it, so that the library tool
// and the app say the same about a file, word for word, as docs/library-format.md promises.
// Changing a check or its words on one side fails the other's tests until both agree.

// sharedCases returns the cases' names, by their .want files.
func sharedCases(t *testing.T) []string {
	t.Helper()
	wants, err := filepath.Glob(filepath.Join("testdata", "checks", "*.want"))
	if err != nil || len(wants) == 0 {
		t.Fatalf("no cases in testdata/checks: %v", err)
	}
	var names []string
	for _, w := range wants {
		names = append(names, strings.TrimSuffix(filepath.Base(w), ".want"))
	}
	return names
}

// sharedCase makes the case's file at path: its .file as it is, or library.sql's library
// changed by its .sql, run with foreign keys off and outside a transaction.
func sharedCase(t *testing.T, name, path string) {
	t.Helper()
	dir := filepath.Join("testdata", "checks")
	if raw, err := os.ReadFile(filepath.Join(dir, name+".file")); err == nil {
		if err := os.WriteFile(path, raw, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	library, err := os.ReadFile(filepath.Join(dir, "library.sql"))
	if err != nil {
		t.Fatal(err)
	}
	change, err := os.ReadFile(filepath.Join(dir, name+".sql"))
	if err != nil {
		t.Fatalf("case %s has neither a .sql nor a .file: %v", name, err)
	}
	d, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	d.SetMaxOpenConns(1) // the pragmas below hold for one connection
	for _, s := range []string{Schema, string(library), `PRAGMA foreign_keys = OFF`, string(change)} {
		if strings.TrimSpace(s) == "" {
			continue
		}
		if _, err := d.Exec(s); err != nil {
			t.Fatalf("making case %s: %v", name, err)
		}
	}
}

func TestSharedChecks(t *testing.T) {
	for _, name := range sharedCases(t) {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "library.sqlite")
			sharedCase(t, name, path)
			// As library check reads it.
			d, _, err := Inspect(path)
			if err != nil {
				t.Fatal(err)
			}
			defer d.Close()
			problems, err := Validate(context.Background(), d)
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join("testdata", "checks", name+".want"))
			if err != nil {
				t.Fatal(err)
			}
			got := ""
			for _, p := range problems {
				got += p + "\n"
			}
			if got != string(want) {
				t.Errorf("the checks report\n%s\nwant (testdata/checks/%s.want)\n%s", got, name, want)
			}
		})
	}
}

package db

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A small reference of an invented package, orbit, whose rows pass every check.
var referenceRows = []string{
	`INSERT INTO source VALUES (1, 'go1.99.0', 'darwin/arm64', '0123456789abcdef0123456789abcdef01234567', '2026-11-01T12:00:00Z', 'By The Go Authors.', 'Copyright 2009 The Go Authors.')`,
	`INSERT INTO package VALUES (1, 'space/orbit', 'orbit', 'Package orbit plans passes.', 'Package orbit plans passes. See [Craft](godoc:space/orbit.Craft) and [the spec](godoc-article:spec#Types).', 0)`,
	`INSERT INTO declaration VALUES (1, 1, 'const', 'const (
	Low  = 1
	High = 2
)', 'Altitudes; see [Craft.Raise](godoc:space/orbit.Craft.Raise).')`,
	`INSERT INTO declaration VALUES (2, 1, 'type', 'type Craft struct{ Name string }', 'A Craft flies. Its [Craft.Name](godoc:space/orbit.Craft.Name) is printed.')`,
	`INSERT INTO declaration VALUES (3, 1, 'method', 'func (c *Craft) Raise()', NULL)`,
	`INSERT INTO symbol VALUES (1, 1, 1, 'const', NULL, 'Low', NULL, 0, 0), (2, 1, 1, 'const', NULL, 'High', NULL, 0, 1)`,
	`INSERT INTO symbol VALUES (3, 1, 2, 'type', NULL, 'Craft', 'A Craft flies.', 0, 2), (4, 1, 3, 'method', 'Craft', 'Raise', NULL, 1, 3)`,
	`INSERT INTO example VALUES (1, 1, NULL, NULL, 0, NULL, 'fmt.Println(1)', '1' || char(10), 0, 0), (2, 1, 3, 'second', 1, 'Raise it.', 'package orbit_test', NULL, 0, 1)`,
	`INSERT INTO article VALUES ('spec', 'The Go Programming Language Specification', 'Language version go1.99', NULL, 0)`,
	`INSERT INTO article_section VALUES ('spec', 0, 'Types', 2, 'Types', 'See [Types](godoc-article:spec#Types) and [orbit](godoc:space/orbit).'), ('spec', 1, 'Boolean_types', 3, 'Boolean types', NULL)`,
	`INSERT INTO search_entry (kind, package_id, title, body) VALUES ('package', 1, 'space/orbit', 'Package orbit plans passes.')`,
	`INSERT INTO search_entry (kind, symbol_id, title, body) SELECT 'symbol', s.id, 'space/orbit.' || coalesce(s.receiver || '.', '') || s.name, s.synopsis FROM symbol s`,
	`INSERT INTO search_entry (kind, article_slug, section_position, title, body) SELECT 'section', article_slug, position, title, NULL FROM article_section`,
}

func newReference(t *testing.T, stmts ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "reference", "go.sqlite")
	d, err := CreateReference(path)
	if err != nil {
		t.Fatal(err)
	}
	d.Close()
	raw(t, path, append([]string{`PRAGMA foreign_keys = ON`}, stmts...)...)
	return path
}

func validateReference(t *testing.T, path string) Problems {
	t.Helper()
	d, err := readOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	p, err := ValidateReference(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCreateReference(t *testing.T) {
	path := newReference(t, referenceRows...)
	d, err := OpenReference(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, tc := range []struct{ query, want string }{
		{`PRAGMA journal_mode`, "delete"},
		{`PRAGMA user_version`, "1"},
		{`PRAGMA application_id`, "1196641618"},
		{`SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name IN ('source', 'package', 'declaration', 'symbol', 'example', 'article', 'article_section', 'search_entry', 'search')`, "9"},
		{`SELECT group_concat(e.title, ' ') FROM search JOIN search_entry e ON e.id = search.rowid WHERE search MATCH '"orbit" "craft" "rai"*'`, "space/orbit.Craft.Raise"},
		{`SELECT group_concat(e.title, ' ') FROM search JOIN search_entry e ON e.id = search.rowid WHERE search MATCH 'boolean'`, "Boolean types"},
	} {
		var got string
		if err := d.QueryRow(tc.query).Scan(&got); err != nil {
			t.Fatalf("%s: %v", tc.query, err)
		}
		if !strings.EqualFold(got, tc.want) {
			t.Errorf("%s = %s, want %s", tc.query, got, tc.want)
		}
	}
	if ReferenceVersion != 1 || ReferenceApplicationID != 0x47534D52 {
		t.Errorf("ReferenceVersion, ReferenceApplicationID = %d, %#x; want what reference.sql sets", ReferenceVersion, ReferenceApplicationID)
	}
	if p, err := ValidateReference(context.Background(), d); err != nil || len(p) > 0 {
		t.Errorf("a whole reference fails its checks: %q %v", p, err)
	}
	if _, err := d.Exec(`DELETE FROM example`); err == nil {
		t.Error("OpenReference opened the file for writing")
	}
	if d, err := CreateReference(path); err == nil {
		d.Close()
		t.Error("CreateReference made a database where one already was")
	}
}

func TestOpenReferenceRefuses(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup []string
		err   string
	}{
		{"a course library", []string{`PRAGMA application_id = 1196641619`}, "is not a Gosmos Go reference (application_id 1196641619, not 1196641618)"},
		{"a newer version", []string{`PRAGMA user_version = 2`}, "is a Go reference of version 2; this program reads version 1 only"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := newReference(t, tc.setup...)
			before, _ := os.ReadFile(path)
			if d, err := OpenReference(path); err == nil || !strings.Contains(err.Error(), tc.err) {
				if err == nil {
					d.Close()
				}
				t.Fatalf("OpenReference: err = %v, want %q", err, tc.err)
			}
			if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
				t.Error("a refused file was changed")
			}
		})
	}
	if _, err := OpenReference(filepath.Join(t.TempDir(), "none.sqlite")); !errors.Is(err, ErrNotFound) {
		t.Errorf("OpenReference of a missing file: err = %v, want ErrNotFound", err)
	}
}

// The schema refuses rows that break its rules.
func TestReferenceSchemaChecks(t *testing.T) {
	path := newReference(t, referenceRows...)
	d, err := connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, tc := range []struct{ name, stmt string }{
		{"a second source row", `INSERT INTO source VALUES (2, 'go1.99.0', 'darwin/arm64', '0123456789abcdef0123456789abcdef01234567', '2026-11-01T12:00:00Z', 'c', 'l')`},
		{"a short website commit", `UPDATE source SET website_commit = 'b8d34bb'`},
		{"a version without go", `UPDATE source SET go_version = '1.27.3'`},
		{"a local generation time", `UPDATE source SET generated_at = '2026-11-01 12:00'`},
		{"an empty doc", `UPDATE package SET doc = ''`},
		{"a method without a receiver", `INSERT INTO symbol VALUES (9, 1, 3, 'method', NULL, 'Lower', NULL, 0, 9)`},
		{"a func with a receiver", `INSERT INTO declaration VALUES (9, 1, 'func', 'func Plan()', NULL); INSERT INTO symbol VALUES (9, 1, 9, 'func', 'Craft', 'Plan', NULL, 0, 9)`},
		{"a second top-level name", `INSERT INTO symbol VALUES (9, 1, 1, 'const', NULL, 'Low', NULL, 0, 9)`},
		{"a second method of one name", `INSERT INTO symbol VALUES (9, 1, 3, 'method', 'Craft', 'Raise', NULL, 0, 9)`},
		{"a symbol of another kind than its declaration", `INSERT INTO symbol VALUES (9, 1, 1, 'var', NULL, 'Mid', NULL, 0, 9)`},
		{"a symbol of another package than its declaration", `INSERT INTO package VALUES (2, 'space/relay', 'relay', NULL, NULL, 0); INSERT INTO symbol VALUES (9, 2, 1, 'const', NULL, 'Mid', NULL, 0, 0)`},
		{"an example of another package's symbol", `INSERT INTO package VALUES (2, 'space/relay', 'relay', NULL, NULL, 0); INSERT INTO example VALUES (9, 2, 3, NULL, 0, NULL, 'x', NULL, 0, 0)`},
		{"unordered output without output", `UPDATE example SET output = NULL, unordered = 1 WHERE id = 1`},
		{"a suffix in capitals", `UPDATE example SET suffix = 'Second' WHERE id = 2`},
		{"an article slug with a space", `INSERT INTO article VALUES ('effective go', 'Effective Go', NULL, NULL, 1)`},
		{"a second section of one anchor", `INSERT INTO article_section VALUES ('spec', 2, 'Types', 2, 'Types again', NULL)`},
		{"a heading level of 7", `UPDATE article_section SET level = 7 WHERE position = 1`},
		{"a search entry for two things", `UPDATE search_entry SET symbol_id = 1 WHERE kind = 'package'`},
		{"a search entry for a missing section", `INSERT INTO search_entry (kind, article_slug, section_position, title) VALUES ('section', 'spec', 9, 'Nowhere')`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := d.Begin()
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()
			if _, err := tx.Exec(tc.stmt); err == nil {
				t.Errorf("%s: accepted", tc.stmt)
			}
		})
	}
}

// Each check fails on a reference broken for it, with a line that names the problem.
func TestValidateReference(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stmts []string
		want  []string
	}{
		{name: "a whole reference"},
		{
			name:  "another kind of file",
			stmts: []string{`PRAGMA application_id = 42`},
			want:  []string{"application_id is 42, not 1196641618: the file is not a Gosmos Go reference"},
		},
		{
			name:  "a newer version",
			stmts: []string{`PRAGMA user_version = 2`},
			want:  []string{"user_version is 2, not 1: the reference is of a version this program does not read"},
		},
		{
			name:  "no source row",
			stmts: []string{`DELETE FROM source`},
			want:  []string{"source holds 0 rows, not exactly 1"},
		},
		{
			name:  "a method of a type that is not there",
			stmts: []string{`UPDATE symbol SET receiver = 'Probe' WHERE id = 4`},
			want:  []string{"method space/orbit.Probe.Raise: its package has no type Probe"},
		},
		{
			name:  "a declaration without symbols",
			stmts: []string{`INSERT INTO declaration VALUES (9, 1, 'var', 'var Drift = 3', NULL)`},
			want:  []string{"declaration 9 of space/orbit declares no symbol"},
		},
		{
			name:  "a section search does not find",
			stmts: []string{`DELETE FROM search_entry WHERE kind = 'section' AND section_position = 1`},
			want:  []string{"search_entry holds 6 rows, want one for each package, symbol and section: 7"},
		},
		{
			name:  "a search index out of step",
			stmts: []string{`INSERT INTO search(search, rowid, title, body) VALUES ('delete', 7, 'Boolean types', NULL)`},
			want:  []string{"the search index holds 6 rows, not search_entry's 7"},
		},
		{
			name: "links to what is not there",
			stmts: []string{
				`UPDATE package SET doc = 'See [Probe](godoc:space/orbit.Probe), [relay](godoc:space/relay), [Low.X](godoc:space/orbit.Low.X) and [go.dev](https://go.dev/).'`,
				`UPDATE article_section SET body = 'See [maps](godoc-article:spec#Map_types) and [mem](godoc-article:mem).' WHERE position = 0`,
			},
			want: []string{
				"package space/orbit links to godoc:space/orbit.Probe, which is not in the reference",
				"package space/orbit links to godoc:space/relay, which is not in the reference",
				"package space/orbit links to godoc:space/orbit.Low.X, which is not in the reference",
				"package space/orbit links to https://go.dev/, which is not in the reference",
				"section Types of article spec links to godoc-article:spec#Map_types, which is not in the reference",
				"section Types of article spec links to godoc-article:mem, which is not in the reference",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := newReference(t, referenceRows...)
			raw(t, path, tc.stmts...)
			before, _ := os.ReadFile(path)
			got := validateReference(t, path)
			if len(got) != len(tc.want) {
				t.Fatalf("problems = %q, want %q", got, tc.want)
			}
			for i, w := range tc.want {
				if got[i] != w {
					t.Errorf("problem %d = %q, want %q", i, got[i], w)
				}
			}
			if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
				t.Error("ValidateReference changed the file")
			}
		})
	}
}

func TestValidateReferenceOfANonDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notes.sqlite")
	os.WriteFile(path, []byte("these are notes, not a database\n"), 0o644)
	got := validateReference(t, path)
	if len(got) != 1 || !strings.HasPrefix(got[0], "the file is not an SQLite database SQLite can read: ") {
		t.Errorf("problems = %q", got)
	}
}

func TestSplitGodoc(t *testing.T) {
	for _, tc := range []struct{ in, path, names string }{
		{"strings", "strings", ""},
		{"strings.Builder", "strings", "Builder"},
		{"io.Reader.Read", "io", "Reader Read"},
		{"math/rand/v2.Int", "math/rand/v2", "Int"},
		{"crypto/sha256", "crypto/sha256", ""},
	} {
		path, names := SplitGodoc(tc.in)
		if path != tc.path || strings.Join(names, " ") != tc.names {
			t.Errorf("SplitGodoc(%q) = %q, %q; want %q, %q", tc.in, path, names, tc.path, tc.names)
		}
	}
}

package db

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// testInfo is the library_info row of the libraries the tests make.
var testInfo = Info{UUID: "6f1d8a52-3c1e-4b7a-9d2e-5a0c7e3b9f14", Title: "Learn Passes", Version: "1"}

func newDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "courses", "passes.sqlite")
	d, err := Create(path, testInfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d, path
}

// raw runs statements on the file without the settings Open and Create make, foreign keys
// included, to build the broken databases the checks must catch.
func raw(t *testing.T, path string, stmts ...string) {
	t.Helper()
	d, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	d.SetMaxOpenConns(1) // so a pragma holds for the statements after it
	for _, s := range stmts {
		if _, err := d.Exec(s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
}

func TestCreate(t *testing.T) {
	d, _ := newDB(t)
	for _, tc := range []struct {
		query string
		want  string
	}{
		{`PRAGMA journal_mode`, "delete"},
		{`PRAGMA foreign_keys`, "1"},
		{`PRAGMA user_version`, "7"},
		{`PRAGMA application_id`, "1196641619"},
		{`SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name IN ('library_info', 'course', 'chapter', 'lesson', 'lesson_topic', 'lesson_hint', 'code_lesson', 'starter_file', 'solution_file', 'quiz_lesson', 'quiz_question', 'quiz_answer', 'cli_lesson', 'cli_step', 'cli_http_test', 'github_lesson', 'interview_lesson')`, "17"},
		{`SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name IN ('progress', 'attempt', 'note', 'code_file', 'project', 'server_address', 'setting')`, "0"},
		{`SELECT count(*) FROM sqlite_schema WHERE type = 'trigger'`, "3"},
		{`SELECT uuid || ' ' || title || ' ' || version FROM library_info`, testInfo.UUID + " Learn Passes 1"},
	} {
		var got string
		if err := d.QueryRow(tc.query).Scan(&got); err != nil {
			t.Fatalf("%s: %v", tc.query, err)
		}
		if !strings.EqualFold(got, tc.want) {
			t.Errorf("%s = %s, want %s", tc.query, got, tc.want)
		}
	}
	if Version != 7 || ApplicationID != 1196641619 || LearnerVersion != 1 || LearnerApplicationID != 1196641612 {
		t.Errorf("Version, ApplicationID, LearnerVersion, LearnerApplicationID = %d, %d, %d, %d; want what the schemas set",
			Version, ApplicationID, LearnerVersion, LearnerApplicationID)
	}
	if p, err := Validate(context.Background(), d); err != nil || len(p) > 0 {
		t.Errorf("a new database fails its checks: %v %v", p, err)
	}
	if info, err := ReadInfo(context.Background(), d); err != nil || info != testInfo {
		t.Errorf("ReadInfo = %+v, %v; want %+v", info, err, testInfo)
	}
}

func TestCreateMakesAUUID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.sqlite")
	d, err := Create(path, Info{Title: "Orbits", Version: "2026.1"})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	info, err := ReadInfo(context.Background(), d)
	if err != nil || !reUUID.MatchString(info.UUID) || info.Title != "Orbits" || info.Version != "2026.1" {
		t.Errorf("ReadInfo = %+v, %v", info, err)
	}
}

func TestCreateRefusesAnEmptyTitle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "untitled.sqlite")
	if d, err := Create(path, Info{Version: "1"}); err == nil {
		d.Close()
		t.Fatal("Create made a library without a title")
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("a refused Create left a file behind")
	}
}

var reUUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestNewUUID(t *testing.T) {
	seen := map[string]bool{}
	for range 1000 {
		u := NewUUID()
		if !reUUID.MatchString(u) || seen[u] {
			t.Fatalf("NewUUID = %q: not a new version 4 uuid", u)
		}
		seen[u] = true
	}
}

func TestCreateRefusesAnExistingPath(t *testing.T) {
	_, path := newDB(t)
	empty := filepath.Join(t.TempDir(), "empty.sqlite")
	os.WriteFile(empty, nil, 0o644)
	for _, p := range []string{path, empty, t.TempDir()} {
		if d, err := Create(p, testInfo); err == nil {
			d.Close()
			t.Errorf("Create(%s) made a database where something already was", p)
		}
	}
}

func TestOpen(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup []string // run on a new database before it is opened again
		ok    bool
		err   string
	}{
		{name: "a database it made", ok: true},
		{name: "one left in WAL mode", setup: []string{`PRAGMA journal_mode = WAL`}, ok: true},
		{name: "a library of format 6", setup: []string{`PRAGMA application_id = 1196447556`, `PRAGMA user_version = 6`}, err: "the library is of format 6, older than 7; library migrate brings it up to date"},
		{name: "an older format", setup: []string{`PRAGMA user_version = 6`}, err: "older than 7"},
		{name: "a newer format", setup: []string{`PRAGMA user_version = 8`}, err: "newer than this program's 7"},
		{name: "not a Gosmos file", setup: []string{`PRAGMA application_id = 42`}, err: "not a Gosmos library (application_id 42)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, path := newDB(t)
			d.Close()
			raw(t, path, tc.setup...)
			before, _ := os.ReadFile(path)
			d, err := Open(path)
			if err == nil {
				defer d.Close()
			}
			if (err == nil) != tc.ok || err != nil && !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("Open: err = %v, want ok = %v, %q", err, tc.ok, tc.err)
			}
			if !tc.ok {
				if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
					t.Error("a refused file was changed")
				}
				return
			}
			var mode string
			d.QueryRow(`PRAGMA journal_mode`).Scan(&mode)
			if mode != "delete" {
				t.Errorf("journal_mode = %s", mode)
			}
		})
	}
}

func TestOpenMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "none.sqlite")
	if _, err := Open(path); !errors.Is(err, ErrNotFound) {
		t.Errorf("Open of a missing file: err = %v, want ErrNotFound", err)
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("Open created the file")
	}
}

func TestOpenWritesNothing(t *testing.T) {
	d, path := newDB(t)
	d.Close()
	before, _ := os.ReadFile(path)
	d, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if p, err := Validate(context.Background(), d); err != nil || len(p) > 0 {
		t.Fatal(p, err)
	}
	d.Close()
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("opening and checking a database changed the file")
	}
	for _, suffix := range []string{"-journal", "-wal", "-shm"} {
		if _, err := os.Stat(path + suffix); err == nil {
			t.Errorf("%s left behind", path+suffix)
		}
	}
}

const (
	course  = `INSERT INTO course VALUES ('c1', 'learn-passes', 'Learn Passes', 'Course', 'Short.', 'Long.', 0)`
	chapter = `INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-intro', 'Intro', '')`
)

func lesson(uuid, typ string, position int) string {
	return `INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod, code_last_mod, synced_at)
		VALUES ('` + uuid + `', 'ch1', ` + string(rune('0'+position)) + `, '` + uuid + `', 'A lesson', '` + typ + `', 1, 10, 'Read me.',
		'2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', NULL, '2026-01-03T00:00:00Z')`
}

// validate checks the file at path as `library check` does: read-only, as it is.
func validate(t *testing.T, path string) Problems {
	t.Helper()
	d, _, err := Inspect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	p, err := Validate(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// corruptPage overwrites the start of a table's root page, as a damaged disk might.
func corruptPage(t *testing.T, path, table string) {
	t.Helper()
	d, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	var root, size int64
	if err := d.QueryRow(`SELECT (SELECT rootpage FROM sqlite_schema WHERE name = ?), (SELECT page_size FROM pragma_page_size)`, table).Scan(&root, &size); err != nil {
		t.Fatal(err)
	}
	d.Close()
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteAt(bytes.Repeat([]byte{0xff}, 16), (root-1)*size); err != nil {
		t.Fatal(err)
	}
}

// Each of the format's rules fails on a file broken for it, with a line that names the problem.
func TestValidate(t *testing.T) {
	quiz := []string{lesson("l1", "type_choice", 0), `INSERT INTO quiz_lesson VALUES ('l1', 'type_choice', 0)`, `INSERT INTO quiz_question VALUES ('l1', 0, 'Why?')`}
	for _, tc := range []struct {
		name   string
		stmts  []string
		damage string // a table whose root page gets overwritten
		want   []string
	}{
		{
			name:  "a whole quiz lesson",
			stmts: append(quiz, `INSERT INTO quiz_answer VALUES ('l1', 0, 0, 'Because.', 1, NULL)`, `INSERT INTO quiz_answer VALUES ('l1', 0, 1, 'No.', 0, NULL)`),
		},
		{
			name:  "another kind of file",
			stmts: []string{`PRAGMA application_id = 42`},
			want:  []string{"application_id is 42, not 1196641619: the file is not a Gosmos library"},
		},
		{
			name:  "a library of format 6",
			stmts: []string{`PRAGMA application_id = 1196447556`, `PRAGMA user_version = 6`},
			want:  []string{"application_id is 1196447556 and user_version 6: the library is of format 6, from before format 7; library migrate brings it up to date"},
		},
		{
			name:  "an older format",
			stmts: []string{`PRAGMA user_version = 6`},
			want:  []string{"user_version is 6, not 7: the library is of an older format; library migrate brings it up to date"},
		},
		{
			name:  "a newer format",
			stmts: []string{`PRAGMA user_version = 8`},
			want:  []string{"user_version is 8, not 7: the library is of a newer format than this program knows"},
		},
		{
			name:  "no library_info row",
			stmts: []string{`DELETE FROM library_info`},
			want:  []string{"library_info holds 0 rows, not exactly 1"},
		},
		{
			name:   "a damaged page",
			stmts:  append(quiz, `INSERT INTO quiz_answer VALUES ('l1', 0, 0, 'Because.', 1, NULL)`),
			damage: "quiz_answer",
			want:   []string{"quick_check: "},
		},
		{
			name:  "not an SQLite file",
			stmts: nil,
			want:  []string{"the file is not an SQLite database SQLite can read: "},
		},
		{
			name:  "a hint for a lesson that is not there",
			stmts: []string{`PRAGMA foreign_keys = OFF`, `INSERT INTO lesson_hint VALUES ('nowhere', 0, 'Look closer.')`},
			want:  []string{"foreign key: lesson_hint row 1 points at a missing lesson row"},
		},
		{
			name:  "a lesson without its per-type row",
			stmts: []string{lesson("l1", "type_choice", 0)},
			want:  []string{"lesson l1 (type_choice) has 0 quiz_lesson rows, want 1"},
		},
		{
			name:  "a lesson with another type's row",
			stmts: []string{lesson("l1", "type_choice", 0), `PRAGMA foreign_keys = OFF`, `INSERT INTO interview_lesson VALUES ('l1', 'type_interview_chat', 'Ask.', 'Answer.')`},
			want: []string{"foreign key: interview_lesson row 1 points at a missing lesson row",
				"lesson l1 (type_choice) has 0 quiz_lesson rows, want 1", "lesson l1 (type_choice) has 1 interview_lesson rows, want 0"},
		},
		{
			name:  "a quiz question without a right answer",
			stmts: append(quiz, `INSERT INTO quiz_answer VALUES ('l1', 0, 0, 'No.', 0, NULL)`),
			want:  []string{"quiz question 0 of lesson l1 has 0 right answers, want 1"},
		},
		{
			name:  "a quiz question without answers",
			stmts: quiz,
			want:  []string{"quiz question 0 of lesson l1 has 0 right answers, want 1"},
		},
		{
			name:  "a command step without its command",
			stmts: []string{lesson("l1", "type_cli", 0), `INSERT INTO cli_lesson VALUES ('l1', 'type_cli', NULL, 0, 0, 0)`, `INSERT INTO cli_step VALUES ('l1', 0, 'command', 'Run it.', 0, NULL)`},
			want:  []string{"cli_step 0 of lesson l1 (command) has no cli_command row"},
		},
		{
			name: "solutions for files the learner does not edit",
			stmts: []string{lesson("l1", "type_code_tests", 0), `INSERT INTO code_lesson VALUES ('l1', 'type_code_tests', 'go', NULL, NULL)`,
				`INSERT INTO starter_file VALUES ('l1', 'sum.go', 0, 'package main', 0, 0), ('l1', 'sum_test.go', 1, 'package main', 0, 1)`,
				`INSERT INTO solution_file VALUES ('l1', 'sum.go', 'package main // solved'), ('l1', 'sum_test.go', 'package main'), ('l1', 'extra.go', 'package main')`},
			want: []string{"solution_file extra.go of lesson l1 solves no editable starter file",
				"solution_file sum_test.go of lesson l1 solves no editable starter file"},
		},
		{
			name:  "an HTTP step without its request",
			stmts: []string{lesson("l1", "type_cli", 0), `INSERT INTO cli_lesson VALUES ('l1', 'type_cli', 'http://localhost:8080', 0, 0, 0)`, `INSERT INTO cli_step VALUES ('l1', 0, 'http', 'Ask the server.', 0, NULL)`},
			want:  []string{"cli_step 0 of lesson l1 (http) has no cli_http_request row"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, path := newDB(t)
			d.Close()
			raw(t, path, append([]string{course, chapter}, tc.stmts...)...)
			if tc.damage != "" {
				corruptPage(t, path, tc.damage)
			}
			if tc.name == "not an SQLite file" {
				os.WriteFile(path, []byte("these are notes, not a database\n"), 0o644)
			}
			before, _ := os.ReadFile(path)
			got := validate(t, path)
			if tc.damage != "" { // as many lines as the damage causes, each from quick_check
				if len(got) == 0 {
					t.Fatal("no problems found in a damaged file")
				}
				for _, g := range got {
					if !strings.HasPrefix(g, "quick_check: ") || strings.Contains(g, "\n") {
						t.Errorf("problem %q is not one quick_check line", g)
					}
				}
				return
			}
			if len(got) != len(tc.want) {
				t.Fatalf("problems = %q, want %d like %q", got, len(tc.want), tc.want)
			}
			for i, w := range tc.want {
				if !strings.HasPrefix(got[i], w) {
					t.Errorf("problem %d = %q, want %q", i, got[i], w)
				}
			}
			if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
				t.Error("Validate changed the file")
			}
		})
	}
}

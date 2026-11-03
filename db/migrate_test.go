package db

// Every library here is invented: courses of a fictional space startup.

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"
)

// schemaOfFormat returns the kept schema of an older format.
func schemaOfFormat(t *testing.T, v int) string {
	t.Helper()
	s, err := formatSchema(v)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func oldLesson(uuid, typ string, position int) string {
	return `INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod, code_last_mod, synced_at)
		VALUES ('` + uuid + `', 'ch1', ` + fmt.Sprint(position) + `, '` + uuid + `', 'Lesson ` + uuid + `', '` + typ + `', 1, 10, 'Read me.',
		'2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', NULL, '2026-01-03T00:00:00Z')`
}

// v3File makes a library of format 3, of the lesson types format 8 keeps, with the learner's
// data: rows Migrate moves, and project folders, which only command-line lessons used and it
// drops.
func v3File(t *testing.T, extra ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "v3.sqlite")
	raw(t, path, append([]string{
		schemaOfFormat(t, 3),
		`INSERT INTO course VALUES ('c1', 'learn-uplinks', 'Learn Uplinks', 'Course', 'Short.', 'Long.', 0)`,
		`INSERT INTO course VALUES ('c2', 'learn-downlinks', 'Learn Downlinks', 'Course', 'Short.', 'Long.', 1)`,
		`INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-link', 'Link', 'Hear the craft.')`,
		oldLesson("ping", "type_code", 0), oldLesson("quiz", "type_choice", 1),
		`INSERT INTO code_lesson VALUES ('ping', 'type_code', 'go', 'pong' || char(10))`,
		`INSERT INTO starter_file VALUES ('ping', 'main.go', 0, 'package main', 0, 0)`,
		`INSERT INTO quiz_lesson VALUES ('quiz', 'type_choice', 0)`,
		`INSERT INTO quiz_question VALUES ('quiz', 0, 'Why?')`,
		`INSERT INTO quiz_answer VALUES ('quiz', 0, 0, 'Because.', 1), ('quiz', 0, 1, 'No.', 0)`,
		`INSERT INTO progress VALUES ('ping', '2026-10-31T12:00:00Z', '2026-10-31T12:30:00Z')`,
		`INSERT INTO code_file VALUES ('ping', 'main.go', 'package main', '2026-10-31T12:00:00Z')`,
		`INSERT INTO project VALUES ('c1', '/Users/me/uplinks', 'https://my-app.example.com')`,
		`INSERT INTO project VALUES ('c2', '/Users/me/downlinks', NULL)`,
		`INSERT INTO setting VALUES ('editor.font_size', '14')`,
	}, extra...)...)
	return path
}

// v4File makes a library of format 4 with the learner's data, a server address among it.
func v4File(t *testing.T, extra ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "v4.sqlite")
	raw(t, path, append([]string{
		schemaOfFormat(t, 4),
		`INSERT INTO course VALUES ('c1', 'learn-uplinks', 'Learn Uplinks', 'Course', 'Short.', 'Long.', 0)`,
		`INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-link', 'Link', 'Hear the craft.')`,
		oldLesson("ping", "type_code", 0), oldLesson("quiz", "type_choice", 1),
		`INSERT INTO code_lesson VALUES ('ping', 'type_code', 'go', 'pong' || char(10))`,
		`INSERT INTO starter_file VALUES ('ping', 'main.go', 0, 'package main', 0, 0)`,
		`INSERT INTO quiz_lesson VALUES ('quiz', 'type_choice', 0)`,
		`INSERT INTO quiz_question VALUES ('quiz', 0, 'Why?')`,
		`INSERT INTO quiz_answer VALUES ('quiz', 0, 0, 'Because.', 1), ('quiz', 0, 1, 'No.', 0)`,
		`INSERT INTO progress VALUES ('quiz', NULL, '2026-10-31T12:30:00Z')`,
		`INSERT INTO note VALUES ('ping', 'Practice this one.', '2026-10-31T12:00:00Z')`,
		`INSERT INTO project VALUES ('c1', NULL, 'https://my-app.example.com')`,
		`INSERT INTO server_address VALUES ('http://localhost:3000', 'http://localhost:3300')`,
	}, extra...)...)
	return path
}

// v5File makes a library of format 5: a code lesson with its starter files, a quiz, and the
// learner's data.
func v5File(t *testing.T, extra ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "v5.sqlite")
	raw(t, path, append([]string{
		schemaOfFormat(t, 5),
		`INSERT INTO course VALUES ('c1', 'learn-orbits', 'Learn Orbits', 'Course', 'Short.', 'Long.', 0)`,
		`INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-launch', 'Launch', 'Lift off.')`,
		oldLesson("burn", "type_code_tests", 0), oldLesson("debrief", "type_choice", 1),
		`INSERT INTO code_lesson VALUES ('burn', 'type_code_tests', 'go', NULL)`,
		`INSERT INTO starter_file VALUES ('burn', 'burn.go', 0, 'package main', 0, 0), ('burn', 'burn_test.go', 1, 'package main', 0, 1)`,
		`INSERT INTO quiz_lesson VALUES ('debrief', 'type_choice', 0)`,
		`INSERT INTO progress VALUES ('burn', '2026-10-31T12:00:00Z', '2026-10-31T12:30:00Z')`,
		`INSERT INTO code_file VALUES ('burn', 'burn.go', 'package main // mine', '2026-10-31T12:10:00Z')`,
		`INSERT INTO server_address VALUES ('http://localhost:8080', 'http://localhost:8181')`,
	}, extra...)...)
	return path
}

// v6Content is a library of format 6 without the learner's data: a code lesson, a code_tests
// lesson with its reference solution, and a quiz, in two courses.
var v6Content = []string{
	`INSERT INTO course VALUES ('c1', 'learn-orbits', 'Learn Orbits', 'Course', 'Short.', 'Long.', 0)`,
	`INSERT INTO course VALUES ('c2', 'learn-landings', 'Learn Landings', 'Guided Project', 'Short.', 'Long.', 1)`,
	`INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-launch', 'Launch', 'Lift off.')`,
	oldLesson("window", "type_code", 0), oldLesson("burn", "type_code_tests", 1), oldLesson("checklist", "type_choice", 2),
	`INSERT INTO lesson_topic VALUES ('window', 'printing'), ('burn', 'functions')`,
	`INSERT INTO code_lesson VALUES ('window', 'type_code', 'go', 'Window opens at T-minus 90 seconds' || char(10))`,
	`INSERT INTO code_lesson VALUES ('burn', 'type_code_tests', 'go', NULL)`,
	`INSERT INTO starter_file VALUES ('window', 'main.go', 0, 'package main', 0, 0)`,
	`INSERT INTO starter_file VALUES ('burn', 'burn.go', 0, 'package main', 0, 0), ('burn', 'burn_test.go', 1, 'package main', 0, 1)`,
	`INSERT INTO solution_file VALUES ('burn', 'burn.go', 'package main // solved')`,
	`INSERT INTO quiz_lesson VALUES ('checklist', 'type_choice', 0)`,
	`INSERT INTO quiz_question VALUES ('checklist', 0, 'What comes before ignition?'), ('checklist', 1, 'Who calls a scrub?')`,
	`INSERT INTO quiz_answer VALUES ('checklist', 0, 0, 'The fuel check.', 1), ('checklist', 0, 1, 'The landing.', 0)`,
	`INSERT INTO quiz_answer VALUES ('checklist', 1, 0, 'Anyone on the console.', 1), ('checklist', 1, 1, 'Nobody.', 0)`,
}

// noteBody is a note of the learner's, to look for in the file once it should be gone.
const noteBody = "Remember the reserve: the burn never takes the tank below ten percent."

// v6Learner is the learner's data in a format 6 library of v6Content. The attempts' ids have
// gaps and are not in the order of their times, so a copy that renumbers them must keep their
// order by id. The attempt of a command-line check (10), the project folder and the server
// addresses only served command-line lessons: Migrate drops them.
var v6Learner = []string{
	`INSERT INTO progress VALUES ('window', '2026-10-30T09:00:00Z', '2026-10-30T09:05:00Z'), ('burn', NULL, '2026-10-31T08:00:00Z'), ('checklist', '2026-10-31T08:30:00Z', '2026-10-31T08:30:00Z')`,
	`INSERT INTO attempt VALUES (5, 'window', 'run', 0, 'Window opens at T-minus ? seconds' || char(10), 312, '2026-10-30T09:01:00Z')`,
	`INSERT INTO attempt VALUES (9, 'window', 'submit', 1, 'Window opens at T-minus 90 seconds' || char(10), 290, '2026-10-30T09:04:00Z')`,
	`INSERT INTO attempt VALUES (10, 'window', 'check', 1, 'every step passed', 80, '2026-10-30T09:06:00Z')`,
	`INSERT INTO attempt VALUES (12, 'checklist', 'answer', 1, 'The fuel check.', NULL, '2026-10-31T08:29:00Z')`,
	`INSERT INTO attempt VALUES (14, 'burn', 'run', 0, '--- FAIL: TestBurn', 1204, '2026-10-31T07:59:00Z')`,
	`INSERT INTO note VALUES ('burn', '` + noteBody + `', '2026-10-31T08:01:00Z')`,
	`INSERT INTO code_file VALUES ('window', 'main.go', 'package main // mine', '2026-10-30T09:04:00Z'), ('burn', 'burn.go', 'package main // half done', '2026-10-31T08:00:00Z')`,
	`INSERT INTO project VALUES ('c2', '/Users/me/landings', NULL)`,
	`INSERT INTO server_address VALUES ('http://localhost:8080', 'http://localhost:8181'), ('http://localhost:3000', 'http://localhost:3300')`,
	`INSERT INTO setting VALUES ('editor.font_size', '14'), ('editor.theme', 'night')`,
}

// v6File makes a library of format 6 from v6Content and the statements given.
func v6File(t *testing.T, name string, stmts ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	raw(t, path, append([]string{schemaOfFormat(t, 6)}, stmts...)...)
	return path
}

// v7Info is the library_info row of v7Content, which format 8 keeps.
var v7Info = Info{UUID: "3c9e7a41-5b2d-4f60-8e1a-7d4b2c6f0a95", Title: "Orbit School", Version: "3"}

// v7Content is a library of format 7 with every kind of row the format 8 keeps: a code lesson
// with hints, a solution and its note, a tests lesson, and a quiz with explanations.
var v7Content = []string{
	`INSERT INTO library_info VALUES (1, '` + v7Info.UUID + `', '` + v7Info.Title + `', '` + v7Info.Version + `')`,
	`INSERT INTO course VALUES ('c1', 'learn-orbits', 'Learn Orbits', 'Course', 'Short.', 'Long.', 0)`,
	`INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-launch', 'Launch', 'Lift off.')`,
	lesson("window", "type_code", 0), lesson("checklist", "type_choice", 2),
	`INSERT INTO lesson (id, uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
		VALUES (23, 'burn', 'ch1', 1, 'burn', 'Burn', 'type_code_tests', 1, 9, 'Burn the engine for six seconds.', '2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z')`,
	`INSERT INTO lesson_topic VALUES ('window', 'printing'), ('burn', 'functions')`,
	`INSERT INTO lesson_hint VALUES ('window', 0, 'Print the seconds.'), ('window', 1, 'Println adds the newline.')`,
	`INSERT INTO code_lesson VALUES ('window', 'type_code', 'go', 'Window opens at T-minus 90 seconds' || char(10), 'Println prints the line whole.')`,
	`INSERT INTO code_lesson VALUES ('burn', 'type_code_tests', 'go', NULL, NULL)`,
	`INSERT INTO starter_file VALUES ('window', 'main.go', 0, 'package main', 0, 0)`,
	`INSERT INTO starter_file VALUES ('burn', 'burn.go', 0, 'package main', 0, 0), ('burn', 'burn_test.go', 1, 'package main', 0, 1)`,
	`INSERT INTO solution_file VALUES ('window', 'main.go', 'package main // solved'), ('burn', 'burn.go', 'package main // solved')`,
	`INSERT INTO quiz_lesson VALUES ('checklist', 'type_choice', 0)`,
	`INSERT INTO quiz_question VALUES ('checklist', 0, 'What comes before ignition?')`,
	`INSERT INTO quiz_answer VALUES ('checklist', 0, 0, 'The fuel check.', 1, 'Nothing burns before the tanks are checked.'), ('checklist', 0, 1, 'The landing.', 0, NULL)`,
}

// v7File makes a library of format 7 from the statements given.
func v7File(t *testing.T, name string, stmts ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	raw(t, path, append([]string{schemaOfFormat(t, 7)}, stmts...)...)
	return path
}

// dump renders every row of every table but the ones named, to see they did not change. A
// table's columns are those cols names for it, or else all of them.
func dump(t *testing.T, d *sql.DB, cols map[string]string, skip ...string) string {
	t.Helper()
	var b strings.Builder
	for _, table := range tableNames(t, d) {
		if strings.HasPrefix(table, "lesson_search") || slices.Contains(skip, table) {
			continue
		}
		cols, ok := cols[table]
		if !ok {
			cols = "*"
		}
		rows, err := d.Query(`SELECT ` + cols + ` FROM ` + table + ` ORDER BY 1, 2`)
		if err != nil {
			t.Fatal(err)
		}
		n, _ := rows.Columns()
		for rows.Next() {
			vals := make([]any, len(n))
			ptrs := make([]any, len(n))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			rows.Scan(ptrs...)
			b.WriteString(table)
			for _, v := range vals {
				b.WriteString("|")
				b.WriteString(strings.TrimSpace(strings.ReplaceAll(sprint(v), "\n", " ")))
			}
			b.WriteString("\n")
		}
		rows.Close()
	}
	return b.String()
}

// oldColumns are the columns format 7 kept as they were, of the tables it changed from format 6:
// lesson lost two columns and has its difficulty on a new scale (TestMigrateRescalesDifficulty),
// and code_lesson and quiz_answer gained one each.
var oldColumns = map[string]string{
	"lesson":      "id, uuid, chapter_uuid, position, slug, title, type, required, readme, created_at, last_mod",
	"code_lesson": "lesson_uuid, lesson_type, prog_lang, expected_output",
	"quiz_answer": "lesson_uuid, question_position, position, answer, is_correct",
}

func sprint(v any) string {
	switch v := v.(type) {
	case nil:
		return "NULL"
	case []byte:
		return string(v)
	default:
		return fmt.Sprint(v)
	}
}

func tableNames(t *testing.T, d *sql.DB) []string {
	t.Helper()
	rows, err := d.Query(`SELECT name FROM sqlite_schema WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var n string
		rows.Scan(&n)
		out = append(out, n)
	}
	return out
}

var reQuotedTable = regexp.MustCompile(`^CREATE TABLE "([^"]+)"`)

// schemaOf returns every object of the database's schema with its SQL, without comments and
// with spaces collapsed. A table rebuilt by ALTER TABLE ... RENAME has its name quoted.
func schemaOf(t *testing.T, d *sql.DB) map[string]string {
	t.Helper()
	rows, err := d.Query(`SELECT type, name, coalesce(sql, '') FROM sqlite_schema`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var typ, name, text string
		rows.Scan(&typ, &name, &text)
		var b strings.Builder
		for _, line := range strings.Split(text, "\n") {
			line, _, _ = strings.Cut(line, "--")
			b.WriteString(line + " ")
		}
		out[typ+" "+name] = reQuotedTable.ReplaceAllString(strings.Join(strings.Fields(b.String()), " "), "CREATE TABLE $1")
	}
	return out
}

// sameSchema fails unless d has the schema of a new library, comments aside.
func sameSchema(t *testing.T, d *sql.DB) {
	t.Helper()
	fresh, err := Create(filepath.Join(t.TempDir(), "new.sqlite"), testInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	got, want := schemaOf(t, d), schemaOf(t, fresh)
	for k, w := range want {
		if got[k] != w {
			t.Errorf("%s:\n got %s\nwant %s", k, got[k], w)
		}
	}
	for k := range got {
		if _, ok := want[k]; !ok {
			t.Errorf("%s is left over", k)
		}
	}
}

// migrated migrates the library at path and checks what every migration to format 9 must
// leave: format 9's marks, the rows of every content table it does not change (changed names
// those an older step changes), the library checks passed, the schema of a new library, and
// neither learner tables nor those of the lesson types format 8 removed. From format 6 or
// older, the library gets a new uuid and no teaching content yet. It returns what Migrate did
// and the library, open.
func migrated(t *testing.T, path string, opt MigrateOptions, changed ...string) (Migration, *sql.DB) {
	t.Helper()
	before, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	var from int
	before.QueryRow(`SELECT user_version FROM pragma_user_version`).Scan(&from)
	cols := map[string]string{}
	if from <= 6 {
		cols = oldColumns
	}
	had := tableNames(t, before)
	kept := dump(t, before, cols, slices.Concat(changed, learnerTableNames(), removedTables)...)
	before.Close()
	m, err := Migrate(context.Background(), path, opt)
	if err != nil {
		t.Fatal(err)
	}
	d, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	if m.From != from || m.To != 9 || !reUUID.MatchString(m.Info.UUID) {
		t.Errorf("Migrate = %+v", m)
	}
	var app, version int
	d.QueryRow(`SELECT (SELECT application_id FROM pragma_application_id), (SELECT user_version FROM pragma_user_version)`).Scan(&app, &version)
	if app != 1196641619 || version != 9 {
		t.Fatalf("application_id %d, user_version %d after migrating; want 1196641619 ('GSMS'), 9", app, version)
	}
	var added []string // tables the library did not have before
	for _, table := range tableNames(t, d) {
		if !slices.Contains(had, table) {
			added = append(added, table)
		}
	}
	if got := dump(t, d, cols, slices.Concat(changed, added)...); got != kept {
		t.Errorf("rows changed:\nbefore\n%s\nafter\n%s", kept, got)
	}
	if info, err := ReadInfo(context.Background(), d); err != nil || info != m.Info {
		t.Errorf("library_info = %+v, %v; want %+v", info, err, m.Info)
	}
	if from <= 6 {
		if m.Info.Version != "1" {
			t.Errorf("library_info.version %q, want 1", m.Info.Version)
		}
		for _, q := range []string{
			`SELECT count(*) FROM code_lesson WHERE solution_note IS NOT NULL`,
			`SELECT count(*) FROM quiz_answer WHERE explanation IS NOT NULL`,
			`SELECT count(*) FROM lesson_hint`,
		} {
			if n := count(t, d, q); n != 0 {
				t.Errorf("%s = %d, want 0", q, n)
			}
		}
	}
	if p, err := Validate(context.Background(), d); err != nil || len(p) > 0 {
		t.Error(p, err)
	}
	sameSchema(t, d)
	return m, d
}

func learnerTableNames() []string {
	var names []string
	for _, t := range learnerTables {
		names = append(names, t.name)
	}
	return names
}

func count(t *testing.T, d *sql.DB, query string, args ...any) int {
	t.Helper()
	var n int
	if err := d.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return n
}

// learnerRowsOf renders a learner database's rows of a library: each table's, but for setting,
// which belongs to no library, all of it.
func learnerRowsOf(t *testing.T, path, library string) map[string][]string {
	t.Helper()
	l, err := OpenLearner(path)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	out := map[string][]string{}
	for table, q := range map[string]string{
		"progress":  `SELECT lesson_uuid, coalesce(done_at, 'NULL'), coalesce(last_opened_at, 'NULL') FROM progress WHERE library_uuid = ?1 ORDER BY 1`,
		"attempt":   `SELECT id, lesson_uuid, kind, passed, output, coalesce(duration_ms, -1), at FROM attempt WHERE library_uuid = ?1 ORDER BY id`,
		"note":      `SELECT lesson_uuid, body, updated_at FROM note WHERE library_uuid = ?1 ORDER BY 1`,
		"code_file": `SELECT lesson_uuid, name, content, updated_at FROM code_file WHERE library_uuid = ?1 ORDER BY 1, 2`,
		"setting":   `SELECT key, value FROM setting WHERE ?1 IS NOT NULL ORDER BY 1`,
	} {
		rows, err := l.Query(q, library)
		if err != nil {
			t.Fatal(err)
		}
		cols, _ := rows.Columns()
		for rows.Next() {
			vals := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			rows.Scan(ptrs...)
			var parts []string
			for _, v := range vals {
				parts = append(parts, strings.ReplaceAll(sprint(v), "\n", `\n`))
			}
			out[table] = append(out[table], strings.Join(parts, "|"))
		}
		rows.Close()
	}
	return out
}

var addedAt = time.Date(2026, 11, 1, 9, 30, 0, 0, time.FixedZone("UTC+4", 4*60*60))

func TestMigrateFrom6MovesTheLearnersData(t *testing.T) {
	path := v6File(t, "orbits.sqlite", append(v6Content, v6Learner...)...)
	learner := filepath.Join(t.TempDir(), "Gosmos", "learner.sqlite")
	m, _ := migrated(t, path, MigrateOptions{Learner: learner, Now: addedAt})
	if m.From != 6 || !m.Created || m.Already || m.Learner != learner || m.Info.Title != "Learn Orbits" {
		t.Errorf("Migrate = %+v", m)
	}
	wantMoved := []TableRows{{"progress", 3}, {"attempt", 4}, {"note", 1}, {"code_file", 2}, {"setting", 2}}
	if !slices.Equal(m.Moved, wantMoved) {
		t.Errorf("moved %v, want %v", m.Moved, wantMoved)
	}
	// What only command-line lessons used goes with the library's learner tables.
	wantDropped := []TableRows{{"attempt", 1}, {"project", 1}, {"server_address", 2}}
	if !slices.Equal(m.Dropped, wantDropped) {
		t.Errorf("dropped %v, want %v", m.Dropped, wantDropped)
	}

	// The learner database was made from learner.sql and holds the library and its rows.
	l, err := OpenLearner(learner)
	if err != nil {
		t.Fatal(err)
	}
	abs, _ := filepath.Abs(path)
	var uuid, title, location, added string
	if err := l.QueryRow(`SELECT uuid, title, location, added_at FROM library`).Scan(&uuid, &title, &location, &added); err != nil {
		t.Fatal(err)
	}
	l.Close()
	if uuid != m.Info.UUID || title != "Learn Orbits" || location != abs || added != "2026-11-01T05:30:00Z" {
		t.Errorf("library row %s %q %s %s", uuid, title, location, added)
	}
	got := learnerRowsOf(t, learner, uuid)
	want := map[string][]string{
		"progress": {"burn|NULL|2026-10-31T08:00:00Z", "checklist|2026-10-31T08:30:00Z|2026-10-31T08:30:00Z", "window|2026-10-30T09:00:00Z|2026-10-30T09:05:00Z"},
		"attempt": { // new ids, in the order of the old ones
			`1|window|run|0|Window opens at T-minus ? seconds\n|312|2026-10-30T09:01:00Z`,
			`2|window|submit|1|Window opens at T-minus 90 seconds\n|290|2026-10-30T09:04:00Z`,
			`3|checklist|answer|1|The fuel check.|-1|2026-10-31T08:29:00Z`,
			`4|burn|run|0|--- FAIL: TestBurn|1204|2026-10-31T07:59:00Z`,
		},
		"note":      {"burn|" + noteBody + "|2026-10-31T08:01:00Z"},
		"code_file": {"burn|burn.go|package main // half done|2026-10-31T08:00:00Z", "window|main.go|package main // mine|2026-10-30T09:04:00Z"},
		"setting":   {"editor.font_size|14", "editor.theme|night"},
	}
	for table, w := range want {
		if !slices.Equal(got[table], w) {
			t.Errorf("%s:\n got %q\nwant %q", table, got[table], w)
		}
	}

	// VACUUM leaves no copy of the learner's data in the library's file.
	if b, _ := os.ReadFile(path); bytes.Contains(b, []byte(noteBody)) || bytes.Contains(b, []byte("half done")) {
		t.Error("the library's file still holds the learner's data")
	}
}

// Format 7's difficulty runs from 1 to 10; the scores of older formats, from about 5 to 70,
// become min(10, max(1, (score + 6) / 7)). Each lesson keeps its id, so the search index still
// finds it, and the columns code_last_mod and synced_at are gone.
func TestMigrateRescalesDifficulty(t *testing.T) {
	scores := map[string]int{"one": 1, "five": 5, "eight": 8, "fifteen": 15, "forty-nine": 49, "fifty": 50, "seventy": 70, "ninety": 90}
	want := map[string]int{"one": 1, "five": 1, "eight": 2, "fifteen": 3, "forty-nine": 7, "fifty": 8, "seventy": 10, "ninety": 10}
	stmts := []string{
		`INSERT INTO course VALUES ('c1', 'learn-orbits', 'Learn Orbits', 'Course', 'Short.', 'Long.', 0)`,
		`INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-launch', 'Launch', 'Lift off.')`,
	}
	position := 0
	for _, slug := range slices.Sorted(maps.Keys(scores)) {
		stmts = append(stmts, fmt.Sprintf(`INSERT INTO lesson (id, uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod, code_last_mod, synced_at)
			VALUES (%d, 'l-%s', 'ch1', %d, '%s', 'Lesson %s', 'type_choice', 1, %d, 'The %s lesson asks about orbit%s.',
			'2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', NULL, '2026-01-03T00:00:00Z')`, 100+position*7, slug, position, slug, slug, scores[slug], slug, slug),
			`INSERT INTO quiz_lesson VALUES ('l-`+slug+`', 'type_choice', 0)`)
		position++
	}
	path := v6File(t, "orbits.sqlite", stmts...)
	_, d := migrated(t, path, MigrateOptions{})
	for slug, w := range want {
		var got, id int
		if err := d.QueryRow(`SELECT difficulty, id FROM lesson WHERE slug = ?`, slug).Scan(&got, &id); err != nil {
			t.Fatal(err)
		}
		if got != w {
			t.Errorf("%s: a score of %d became %d, want %d", slug, scores[slug], got, w)
		}
		var found int
		if err := d.QueryRow(`SELECT rowid FROM lesson_search WHERE lesson_search MATCH ?`, `"orbit`+slug+`"`).Scan(&found); err != nil || found != id {
			t.Errorf("%s: the search index finds row %d (%v), want the lesson's id %d", slug, found, err, id)
		}
	}
	for _, column := range []string{"code_last_mod", "synced_at"} {
		if n := count(t, d, `SELECT count(*) FROM pragma_table_info('lesson') WHERE name = ?`, column); n != 0 {
			t.Errorf("lesson still has %s", column)
		}
	}
	for _, tc := range []struct {
		difficulty int
		ok         bool
	}{{0, false}, {1, true}, {10, true}, {11, false}} {
		_, err := d.Exec(`UPDATE lesson SET difficulty = ? WHERE slug = 'one'`, tc.difficulty)
		if (err == nil) != tc.ok {
			t.Errorf("difficulty %d: err = %v, want ok = %v", tc.difficulty, err, tc.ok)
		}
	}
}

// A learner database that already has settings and attempts keeps its own values where the
// library's clash, and the library's attempts follow its own.
func TestMigrateKeepsTheLearnersValues(t *testing.T) {
	learner := filepath.Join(t.TempDir(), "learner.sqlite")
	l, err := CreateLearner(learner)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{
		`INSERT INTO library VALUES ('apogee-uuid', 'Apogee', 'builtin:apogee.sqlite', '2026-10-31T10:00:00Z')`,
		`INSERT INTO attempt (library_uuid, lesson_uuid, kind, passed, output, duration_ms, at) VALUES ('apogee-uuid', 'l1', 'run', 1, 'ok', 10, '2026-10-31T10:05:00Z')`,
		`INSERT INTO setting VALUES ('editor.font_size', '16')`,
	} {
		if _, err := l.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	l.Close()
	path := v6File(t, "orbits.sqlite", append(v6Content, v6Learner...)...)
	m, _ := migrated(t, path, MigrateOptions{Learner: learner, Title: "Orbit School"})
	if m.Created || m.Info.Title != "Orbit School" {
		t.Errorf("Migrate = %+v", m)
	}
	got := learnerRowsOf(t, learner, m.Info.UUID)
	if want := []string{"editor.font_size|16", "editor.theme|night"}; !slices.Equal(got["setting"], want) {
		t.Errorf("setting: %q, want %q", got["setting"], want)
	}
	if len(got["attempt"]) != 4 || !strings.HasPrefix(got["attempt"][0], "2|window|run|") || !strings.HasPrefix(got["attempt"][3], "5|burn|run|") {
		t.Errorf("attempts %q, want ids 2 to 5 in the old order", got["attempt"])
	}
}

func TestMigrateWithoutLearnerData(t *testing.T) {
	for _, tc := range []struct {
		name, file string
		stmts      []string
		title      string // given
		want       string // library_info.title
	}{
		{"the first course's title", "orbits.sqlite", v6Content, "", "Learn Orbits"},
		{"a title given", "orbits.sqlite", v6Content, "  Orbit School ", "Orbit School"},
		{"the file's name", "apogee.sqlite", nil, "", "apogee"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := v6File(t, tc.file, tc.stmts...)
			m, _ := migrated(t, path, MigrateOptions{Title: tc.title})
			if m.Info.Title != tc.want || m.Learner != "" || len(m.Moved) != 0 {
				t.Errorf("Migrate = %+v, want the title %q and no learner database", m, tc.want)
			}
		})
	}
}

func TestMigrateRefuses(t *testing.T) {
	text := filepath.Join(t.TempDir(), "notes.sqlite")
	os.WriteFile(text, []byte("these are notes, not a database\n"), 0o644)
	newer := v6File(t, "newer.sqlite", `PRAGMA application_id = 1196641619`, `PRAGMA user_version = 10`)
	library9 := filepath.Join(t.TempDir(), "other.sqlite")
	if d, err := Create(library9, testInfo); err != nil {
		t.Fatal(err)
	} else {
		d.Close()
	}
	for _, tc := range []struct {
		name    string
		path    string
		learner string
		want    string
	}{
		{"learner data without a learner database", v6File(t, "orbits.sqlite", append(v6Content, v6Learner...)...), "",
			"holds the learner's data (3 progress, 4 attempt, 1 note, 2 code_file and 2 setting rows); give a learner database to move it into (--learner <file>): nothing changed"},
		{"a library as the learner database", v6File(t, "orbits.sqlite", append(v6Content, v6Learner...)...), library9,
			"is not a Gosmos learner database (application_id 1196641619, not 1196641612): nothing changed"},
		{"a text file as the learner database", v6File(t, "orbits.sqlite", append(v6Content, v6Learner...)...), text, "nothing changed"},
		{"a newer format", newer, "", "the library is of format 10, newer than this program's 9"},
		{"a format 7 library without its library_info row", v7File(t, "nameless.sqlite", v7Content[1:]...), "", "the library has no library_info row; nothing changed"},
		{"an older format than any migration takes", v6File(t, "v2.sqlite", `PRAGMA user_version = 2`), "", "is of format 2; library migrate takes format 3 and later only"},
		{"not a library", v6File(t, "other.sqlite", `PRAGMA application_id = 42`), "", "not a Gosmos library (application_id 42)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before, _ := os.ReadFile(tc.path)
			var learnerBefore []byte
			if tc.learner != "" {
				learnerBefore, _ = os.ReadFile(tc.learner)
			}
			_, err := Migrate(context.Background(), tc.path, MigrateOptions{Learner: tc.learner})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Migrate: %v, want %q", err, tc.want)
			}
			if after, _ := os.ReadFile(tc.path); !bytes.Equal(before, after) {
				t.Error("a refused migration changed the library")
			}
			if tc.learner != "" {
				if after, _ := os.ReadFile(tc.learner); !bytes.Equal(learnerBefore, after) {
					t.Error("a refused migration changed the learner database's file")
				}
			}
		})
	}
}

// An interrupted migration runs again safely: the learner database committed first keeps the
// library row and its rows, the rerun finds them by the file's location and copies nothing
// twice, and once the library is of format 9 a rerun does nothing at all.
func TestMigrateRunsAgain(t *testing.T) {
	// A row whose foreign key points nowhere makes the library's step fail after the learner
	// database has committed, as a crash between the two would leave things.
	broken := `INSERT INTO lesson_topic VALUES ('nowhere', 'lost')`
	path := v6File(t, "orbits.sqlite", append(append(v6Content, v6Learner...), `PRAGMA foreign_keys = OFF`, broken)...)
	learner := filepath.Join(t.TempDir(), "learner.sqlite")
	before, _ := os.ReadFile(path)
	if _, err := Migrate(context.Background(), path, MigrateOptions{Learner: learner}); err == nil || !strings.Contains(err.Error(), "foreign keys pointing nowhere") {
		t.Fatalf("Migrate of a broken library: %v", err)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("the failed step changed the library")
	}
	l, err := OpenLearner(learner)
	if err != nil {
		t.Fatal(err)
	}
	var first string
	if err := l.QueryRow(`SELECT uuid FROM library`).Scan(&first); err != nil {
		t.Fatal(err)
	}
	l.Close()
	copied := learnerRowsOf(t, learner, first)

	// The counts must match: with a new attempt in the library, the rerun refuses.
	raw(t, path, `INSERT INTO attempt VALUES (20, 'window', 'run', 1, 'again', 100, '2026-11-01T08:00:00Z')`)
	if _, err := Migrate(context.Background(), path, MigrateOptions{Learner: learner}); err == nil ||
		!strings.Contains(err.Error(), "the learner database already has this library (uuid "+first) ||
		!strings.Contains(err.Error(), "attempt: 5 rows in the file, 4 in the learner database") {
		t.Fatalf("a rerun with other counts: %v", err)
	}
	raw(t, path, `DELETE FROM attempt WHERE id = 20`, `PRAGMA foreign_keys = OFF`, `DELETE FROM lesson_topic WHERE lesson_uuid = 'nowhere'`)

	m, _ := migrated(t, path, MigrateOptions{Learner: learner, Title: "Ignored"})
	if !m.Already || m.Created || m.Info.UUID != first || m.Info.Title != "Learn Orbits" {
		t.Errorf("the rerun: %+v, want the first run's library %s, nothing copied", m, first)
	}
	if again := learnerRowsOf(t, learner, first); fmt.Sprint(again) != fmt.Sprint(copied) {
		t.Errorf("the rerun changed the learner's rows:\n%v\nwant\n%v", again, copied)
	}

	// Once the library is of format 9, migrating it again changes nothing.
	lib, _ := os.ReadFile(path)
	learned, _ := os.ReadFile(learner)
	m, err = Migrate(context.Background(), path, MigrateOptions{Learner: learner})
	if err != nil || m.From != 9 || m.Info.UUID != first {
		t.Errorf("migrating a library of format 9: %+v, %v", m, err)
	}
	if a, _ := os.ReadFile(path); !bytes.Equal(a, lib) {
		t.Error("migrating a library of format 9 changed it")
	}
	if a, _ := os.ReadFile(learner); !bytes.Equal(a, learned) {
		t.Error("migrating a library of format 9 changed the learner database")
	}
}

// Formats 3, 4 and 5 go to 9 one step at a time, each in its own transaction. The learner's rows
// move into the learner database, but for the project folders and server addresses, which only
// command-line lessons used.
func TestMigrateFrom3(t *testing.T) {
	learner := filepath.Join(t.TempDir(), "learner.sqlite")
	m, _ := migrated(t, v3File(t), MigrateOptions{Learner: learner})
	got := learnerRowsOf(t, learner, m.Info.UUID)
	for table, want := range map[string][]string{
		"progress":  {"ping|2026-10-31T12:00:00Z|2026-10-31T12:30:00Z"},
		"code_file": {"ping|main.go|package main|2026-10-31T12:00:00Z"},
		"setting":   {"editor.font_size|14"},
	} {
		if !slices.Equal(got[table], want) {
			t.Errorf("%s: %q, want %q", table, got[table], want)
		}
	}
	if want := []TableRows{{"attempt", 0}, {"project", 2}, {"server_address", 0}}; !slices.Equal(m.Dropped, want) {
		t.Errorf("dropped %v, want %v", m.Dropped, want)
	}
}

func TestMigrateFrom4(t *testing.T) {
	learner := filepath.Join(t.TempDir(), "learner.sqlite")
	m, _ := migrated(t, v4File(t), MigrateOptions{Learner: learner})
	got := learnerRowsOf(t, learner, m.Info.UUID)
	if !slices.Equal(got["progress"], []string{"quiz|NULL|2026-10-31T12:30:00Z"}) || !slices.Equal(got["note"], []string{"ping|Practice this one.|2026-10-31T12:00:00Z"}) {
		t.Errorf("learner rows %v", got)
	}
	if want := []TableRows{{"attempt", 0}, {"project", 1}, {"server_address", 1}}; !slices.Equal(m.Dropped, want) {
		t.Errorf("dropped %v, want %v", m.Dropped, want)
	}
}

func TestMigrateFrom5(t *testing.T) {
	learner := filepath.Join(t.TempDir(), "learner.sqlite")
	m, d := migrated(t, v5File(t), MigrateOptions{Learner: learner})
	for _, tc := range []struct {
		stmt string
		ok   bool
	}{
		{`INSERT INTO solution_file VALUES ('burn', 'burn.go', 'package main // solved')`, true},
		{`INSERT INTO solution_file VALUES ('burn', 'burn.go', 'again')`, false},
		{`INSERT INTO lesson_hint VALUES ('burn', 0, 'Subtract the burn from the tank.')`, true},
		{`INSERT INTO lesson_hint VALUES ('burn', 1, '')`, false},
		{`UPDATE code_lesson SET solution_note = 'It subtracts.' WHERE lesson_uuid = 'burn'`, true},
		{`UPDATE code_lesson SET solution_note = '' WHERE lesson_uuid = 'burn'`, false},
		{`UPDATE library_info SET title = ''`, false},
		{`INSERT INTO library_info VALUES (2, 'another', 'Another', '1')`, false},
	} {
		if _, err := d.Exec(tc.stmt); (err == nil) != tc.ok {
			t.Errorf("%s: err = %v, want ok = %v", tc.stmt, err, tc.ok)
		}
	}
	if p, err := Validate(context.Background(), d); err != nil || len(p) > 0 {
		t.Errorf("format 7's new rows: %v %v", p, err)
	}
	if got := learnerRowsOf(t, learner, m.Info.UUID); !slices.Equal(got["code_file"], []string{"burn|burn.go|package main // mine|2026-10-31T12:10:00Z"}) {
		t.Errorf("code files %q", got["code_file"])
	}
	if want := []TableRows{{"attempt", 0}, {"project", 0}, {"server_address", 1}}; !slices.Equal(m.Dropped, want) {
		t.Errorf("dropped %v, want %v", m.Dropped, want)
	}
}

// Format 7 goes to 9 keeping every row of the lesson types format 8 keeps, its teaching content
// and its library_info row, so the learner's progress, filed under its uuid, stays its own. The
// tables of the removed types go, and lesson.type takes no other type again.
func TestMigrateFrom7(t *testing.T) {
	path := v7File(t, "orbits.sqlite", v7Content...)
	before := filepath.Join(t.TempDir(), "before.sqlite")
	b, _ := os.ReadFile(path)
	os.WriteFile(before, b, 0o644)
	m, d := migrated(t, path, MigrateOptions{Learner: filepath.Join(t.TempDir(), "unused.sqlite"), Title: "Ignored"})
	if m.Info != v7Info || m.Learner != "" || m.Created || len(m.Moved) != 0 || len(m.Dropped) != 0 {
		t.Errorf("Migrate = %+v, want the library's own library_info and no learner data", m)
	}
	old, err := sql.Open("sqlite", before)
	if err != nil {
		t.Fatal(err)
	}
	defer old.Close()
	if got, want := dump(t, d, nil), dump(t, old, nil, removedTables...); got != want {
		t.Errorf("rows changed:\nbefore\n%s\nafter\n%s", want, got)
	}
	var id, found int
	d.QueryRow(`SELECT id FROM lesson WHERE uuid = 'burn'`).Scan(&id)
	if err := d.QueryRow(`SELECT rowid FROM lesson_search WHERE lesson_search MATCH 'burn'`).Scan(&found); err != nil || found != id || id != 23 {
		t.Errorf("the search index finds row %d (%v) for the lesson with id %d, want 23", found, err, id)
	}
	for _, typ := range removedTypes {
		if _, err := d.Exec(`UPDATE lesson SET type = ? WHERE uuid = 'window'`, typ); err == nil {
			t.Errorf("lesson.type took %s", typ)
		}
	}
	for _, table := range removedTables {
		if n := count(t, d, `SELECT count(*) FROM sqlite_schema WHERE name = ?`, table); n != 0 {
			t.Errorf("%s is still there", table)
		}
	}
	if n := count(t, d, `SELECT freelist_count FROM pragma_freelist_count`); n != 0 {
		t.Errorf("%d free pages after migrating: the file was not vacuumed", n)
	}
}

// v8File makes a library of format 8 from the statements given. The rows of v7Content are rows
// of format 8 too: format 8 only took away what the lesson types it removed used.
func v8File(t *testing.T, name string, stmts ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	raw(t, path, append([]string{schemaOfFormat(t, 8)}, stmts...)...)
	return path
}

// Format 8 goes to 9 with every row as it was, its library_info row included, and the tables of
// a course's story new and empty. Their keys and checks hold from the start, and migrating the
// library again changes nothing.
func TestMigrateFrom8(t *testing.T) {
	path := v8File(t, "orbits.sqlite", v7Content...)
	before := filepath.Join(t.TempDir(), "before.sqlite")
	b, _ := os.ReadFile(path)
	os.WriteFile(before, b, 0o644)
	m, d := migrated(t, path, MigrateOptions{Learner: filepath.Join(t.TempDir(), "unused.sqlite"), Title: "Ignored"})
	if m.From != 8 || m.Info != v7Info || m.Learner != "" || m.Created || len(m.Moved) != 0 || len(m.Dropped) != 0 {
		t.Errorf("Migrate = %+v, want the library's own library_info and no learner data", m)
	}
	old, err := sql.Open("sqlite", before)
	if err != nil {
		t.Fatal(err)
	}
	defer old.Close()
	if got, want := dump(t, d, nil), dump(t, old, nil); got != want {
		t.Errorf("rows changed:\nbefore\n%s\nafter\n%s", want, got)
	}
	for _, table := range storyTables {
		if n := count(t, d, `SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name = ?`, table); n != 1 {
			t.Errorf("%s is not there", table)
		} else if n := count(t, d, `SELECT count(*) FROM `+table); n != 0 {
			t.Errorf("%s holds %d rows, want none", table, n)
		}
	}
	var id, found int
	d.QueryRow(`SELECT id FROM lesson WHERE uuid = 'burn'`).Scan(&id)
	if err := d.QueryRow(`SELECT rowid FROM lesson_search WHERE lesson_search MATCH 'burn'`).Scan(&found); err != nil || found != id || id != 23 {
		t.Errorf("the search index finds row %d (%v) for the lesson with id %d, want 23", found, err, id)
	}
	for _, tc := range []struct {
		stmt string
		ok   bool
	}{
		{`INSERT INTO story_page VALUES ('c1', 0, 'badge', 'APOGEE', 'Your Badge', 'Everyone on the console wears one.', 'Print my badge')`, true},
		{`INSERT INTO story_page VALUES ('c1', 1, 'video', NULL, 'The Launch', 'A film of the launch.', 'Play it')`, false},
		{`INSERT INTO story_page VALUES ('c1', 1, 'text', '', 'Meet the Team', 'Four people run the console.', 'Meet them')`, false},
		{`INSERT INTO story_card VALUES ('c1', 0, 0, 'Role', 'Flight Software Engineer', NULL, NULL)`, true},
		{`INSERT INTO story_card VALUES ('c1', 5, 0, 'Role', 'A card of no page', NULL, NULL)`, false},
		{`INSERT INTO chapter_story VALUES ('ch1', 'The countdown starts today.', 'The window opened on time.')`, true},
		{`INSERT INTO chapter_story VALUES ('no-such-chapter', 'Lift off.', 'In orbit.')`, false},
		{`INSERT INTO story_term VALUES ('c1', 'burn', 'Firing the engine for a set time.', 1)`, true},
		{`INSERT INTO story_term VALUES ('c1', 'burn', 'Firing the engine again.', 1)`, false},
		{`INSERT INTO story_term VALUES ('c1', 'scrub', 'Calling off a launch.', 2)`, false},
	} {
		if _, err := d.Exec(tc.stmt); (err == nil) != tc.ok {
			t.Errorf("%s: err = %v, want ok = %v", tc.stmt, err, tc.ok)
		}
	}
	if p, err := Validate(context.Background(), d); err != nil || len(p) > 0 {
		t.Errorf("format 9's new rows: %v %v", p, err)
	}

	lib, _ := os.ReadFile(path)
	if m, err := Migrate(context.Background(), path, MigrateOptions{}); err != nil || m.From != 9 || m.Info != v7Info {
		t.Errorf("migrating a library of format 9: %+v, %v", m, err)
	}
	if a, _ := os.ReadFile(path); !bytes.Equal(a, lib) {
		t.Error("migrating a library of format 9 changed it")
	}
}

// A library that holds a lesson of a type format 8 removed is refused, from every format
// Migrate takes, before anything changes: no step runs and no learner database is made.
func TestMigrateRefusesRemovedTypes(t *testing.T) {
	course := []string{
		`INSERT INTO course VALUES ('c1', 'learn-uplinks', 'Learn Uplinks', 'Course', 'Short.', 'Long.', 0)`,
		`INSERT INTO chapter VALUES ('ch1', 'c1', 0, '1-link', 'Link', 'Hear the craft.')`,
		`PRAGMA foreign_keys = OFF`,
	}
	for v := 3; v <= 7; v++ {
		t.Run(fmt.Sprintf("format %d", v), func(t *testing.T) {
			stmts := append([]string{schemaOfFormat(t, v)}, course...)
			for i, typ := range []string{"type_cli", "type_cli", "type_github_checks", "type_interview_chat", "type_code"} {
				if v == 7 {
					stmts = append(stmts, lesson(fmt.Sprintf("l%d", i), typ, i))
				} else {
					stmts = append(stmts, oldLesson(fmt.Sprintf("l%d", i), typ, i))
				}
			}
			if v <= 6 {
				stmts = append(stmts, `INSERT INTO progress VALUES ('l0', '2026-10-31T12:00:00Z', NULL)`)
			}
			path := filepath.Join(t.TempDir(), "uplinks.sqlite")
			raw(t, path, stmts...)
			learner := filepath.Join(t.TempDir(), "learner.sqlite")
			before, _ := os.ReadFile(path)
			_, err := Migrate(context.Background(), path, MigrateOptions{Learner: learner})
			want := "holds lessons of types format 8 removed (type_cli 2, type_github_checks 1, type_interview_chat 1): Gosmos is self-contained, with code, tests and quiz lessons only, which it checks on the Mac alone; nothing changed"
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("Migrate: %v, want %q", err, want)
			}
			if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
				t.Error("the refused migration changed the library")
			}
			if _, err := os.Stat(learner); err == nil {
				t.Error("the refused migration made a learner database")
			}
		})
	}
}

// A step that would break a foreign key changes nothing, from any format, and makes no
// learner database.
func TestMigrationRollsBack(t *testing.T) {
	for name, path := range map[string]string{
		"format 3": v3File(t, `PRAGMA foreign_keys = OFF`, `INSERT INTO project VALUES ('no-such-course', '/Users/me/x', NULL)`),
		"format 4": v4File(t, `PRAGMA foreign_keys = OFF`, `INSERT INTO note VALUES ('no-such-lesson', 'Lost.', '2026-10-31T12:00:00Z')`),
		"format 5": v5File(t, `PRAGMA foreign_keys = OFF`, `INSERT INTO attempt (lesson_uuid, kind, passed, output, at) VALUES ('no-such-lesson', 'run', 1, '', '2026-10-31T12:00:00Z')`),
	} {
		learner := filepath.Join(t.TempDir(), "learner.sqlite")
		before, _ := os.ReadFile(path)
		if _, err := Migrate(context.Background(), path, MigrateOptions{Learner: learner}); err == nil || !strings.Contains(err.Error(), "foreign keys pointing nowhere") {
			t.Errorf("%s: err = %v", name, err)
		}
		if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
			t.Errorf("%s: the failed migration changed the file", name)
		}
		if _, err := os.Stat(learner); err == nil {
			t.Errorf("%s: the failed migration made a learner database", name)
		}
	}
}

func TestCreateLearner(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Application Support", "Gosmos", "learner.sqlite")
	l, err := CreateLearner(path)
	if err != nil {
		t.Fatal(err)
	}
	var app, version, tables int
	l.QueryRow(`SELECT (SELECT application_id FROM pragma_application_id), (SELECT user_version FROM pragma_user_version),
		(SELECT count(*) FROM sqlite_schema WHERE type = 'table')`).Scan(&app, &version, &tables)
	l.Close()
	if app != 1196641612 || version != 2 || tables != 10 {
		t.Errorf("application_id %d, user_version %d, %d tables; want 1196641612 ('GSML'), 2, 10", app, version, tables)
	}
	if l, err := CreateLearner(path); err == nil {
		l.Close()
		t.Error("CreateLearner made a learner database where one already was")
	}
	if l, err := OpenLearner(path); err != nil {
		t.Error(err)
	} else {
		l.Close()
	}
	raw(t, path, `PRAGMA user_version = 3`)
	if l, err := OpenLearner(path); err == nil || !strings.Contains(err.Error(), "a learner database of version 3; this program writes version 2 only") {
		t.Errorf("OpenLearner of a newer learner database: %v", err)
		if l != nil {
			l.Close()
		}
	}
	raw(t, path, `PRAGMA user_version = 1`)
	if l, err := OpenLearner(path); err == nil || !strings.Contains(err.Error(), "a learner database of version 1, older than 2: Gosmos brings it up to date when it opens it") {
		t.Errorf("OpenLearner of an older learner database: %v", err)
		if l != nil {
			l.Close()
		}
	}
}

func TestFormatSchemas(t *testing.T) {
	for v := 3; v <= 9; v++ {
		s, err := formatSchema(v)
		if err != nil || !regexp.MustCompile(fmt.Sprintf(`(?m)^PRAGMA user_version = %d;`, v)).MatchString(s) {
			t.Errorf("formatSchema(%d) does not set user_version %d: %v", v, v, err)
		}
	}
	if _, err := formatSchema(2); err == nil {
		t.Error("formatSchema(2) should fail")
	}
}

func TestTableDefinition(t *testing.T) {
	for _, tc := range []struct {
		format int
		table  string
	}{{4, "cli_lesson"}, {4, "project"}, {4, "server_address"}, {5, "interview_lesson"}, {6, "solution_file"}, {7, "library_info"}, {7, "lesson_hint"}, {7, "code_lesson"}, {8, "lesson"},
		{9, "story_page"}, {9, "story_card"}, {9, "chapter_story"}, {9, "story_term"}} {
		def, err := tableDefinition(tc.format, tc.table)
		if err != nil || !strings.HasPrefix(def, "CREATE TABLE "+tc.table+" (") || !strings.HasSuffix(strings.TrimSpace(def), "STRICT") {
			t.Errorf("tableDefinition(%d, %s) = %q, %v", tc.format, tc.table, def, err)
		}
	}
	if _, err := tableDefinition(7, "progress"); err == nil {
		t.Error("format 7 has no progress table")
	}
	if _, err := tableDefinition(8, "cli_lesson"); err == nil {
		t.Error("format 8 has no cli_lesson table")
	}
	if _, err := tableDefinition(8, "story_page"); err == nil {
		t.Error("format 8 has no story_page table")
	}
}

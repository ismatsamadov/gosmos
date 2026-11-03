package main

// Every library in these tests is invented here: a small course about testing in Go, with
// lessons of every type.

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	osexec "os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jobloom/gosmos/db"
)

func runCmd(t *testing.T, args ...string) (code int, out, errOut string) {
	t.Helper()
	var o, e bytes.Buffer
	code = run(context.Background(), args, &o, &e)
	return code, o.String(), e.String()
}

// exec runs statements on a file directly, without the checks the tool makes, to build the
// libraries the tests need.
func exec(t *testing.T, path string, stmts ...string) {
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

func lesson(uuid, chapter, typ string, position int) string {
	return `INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
		VALUES ('` + uuid + `', '` + chapter + `', ` + string(rune('0'+position)) + `, '` + uuid + `', 'Lesson ` + uuid + `', '` + typ + `', 1, 3,
		'# A lesson\n\nRead this first.', '2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z')`
}

// smallLibrary is an invented course with a lesson of every type.
var smallLibrary = []string{
	`PRAGMA foreign_keys = ON`,
	`INSERT INTO course VALUES ('c-test', 'learn-testing', 'Learn Testing', 'Course', 'Tests, from the start.', 'A small invented course.', 0)`,
	`INSERT INTO chapter VALUES ('ch-basics', 'c-test', 0, '1-basics', 'Basics', 'First steps.')`,
	`INSERT INTO chapter VALUES ('ch-more', 'c-test', 1, '2-more', 'More', 'Tables of tests.')`,
	lesson("l-hello", "ch-basics", "type_code", 0),
	`INSERT INTO code_lesson VALUES ('l-hello', 'type_code', 'go', 'Hello, tests!' || char(10), 'Println ends the line for you.')`,
	`INSERT INTO lesson_hint VALUES ('l-hello', 0, 'Start from what main prints now.'), ('l-hello', 1, 'fmt.Println adds the newline.')`,
	`INSERT INTO starter_file VALUES ('l-hello', 'main.go', 0, 'package main', 0, 0)`,
	lesson("l-sum", "ch-basics", "type_code_tests", 1),
	`INSERT INTO code_lesson VALUES ('l-sum', 'type_code_tests', 'go', NULL, 'Adding the two numbers is all the test asks for.')`,
	`INSERT INTO lesson_hint VALUES ('l-sum', 0, 'Read what the test expects.'), ('l-sum', 1, 'Return the sum of a and b.'), ('l-sum', 2, 'Use +.')`,
	`INSERT INTO starter_file VALUES ('l-sum', 'sum.go', 0, 'package main', 0, 0), ('l-sum', 'sum_test.go', 1, 'package main', 0, 1)`,
	lesson("l-quiz", "ch-basics", "type_choice", 2),
	`INSERT INTO quiz_lesson VALUES ('l-quiz', 'type_choice', 0)`,
	`INSERT INTO quiz_question VALUES ('l-quiz', 0, 'Which functions does go test run?')`,
	`INSERT INTO quiz_answer VALUES ('l-quiz', 0, 0, 'Functions named TestXxx', 1, 'go test runs the functions whose names start with Test.'), ('l-quiz', 0, 1, 'Every function', 0, 'A helper without the Test prefix never runs on its own.')`,
	lesson("l-table", "ch-more", "type_code_tests", 0),
	`INSERT INTO code_lesson VALUES ('l-table', 'type_code_tests', 'go', NULL, 'Each case is a row of the table.')`,
	`INSERT INTO starter_file VALUES ('l-table', 'double.go', 0, 'package main', 0, 0), ('l-table', 'double_test.go', 1, 'package main', 0, 1)`,
}

func newLibrary(t *testing.T, stmts ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "library.sqlite")
	if code, _, errOut := runCmd(t, "create", path); code != 0 {
		t.Fatalf("create: %d %s", code, errOut)
	}
	exec(t, path, stmts...)
	return path
}

func TestCreateThenCheck(t *testing.T) {
	path := filepath.Join(t.TempDir(), "libraries", "empty.sqlite")
	code, out, errOut := runCmd(t, "create", path)
	if code != 0 || !regexp.MustCompile(`^created `+regexp.QuoteMeta(path)+`: an empty library of format 9, "empty" version 1 \(uuid [0-9a-f-]{36}\)\n$`).MatchString(out) {
		t.Fatalf("create: %d %q %q", code, out, errOut)
	}
	code, out, _ = runCmd(t, "check", path)
	if code != 0 || !regexp.MustCompile(`^ok: `+regexp.QuoteMeta(path)+` is "empty" version 1 \(uuid [0-9a-f-]{36}\), a library of format 9 with 0 courses, 0 chapters and 0 lessons\n$`).MatchString(out) {
		t.Errorf("check of an empty library: %d %q", code, out)
	}

	titled := filepath.Join(t.TempDir(), "apogee.sqlite")
	if code, out, errOut := runCmd(t, "create", "--title", "Apogee", titled, "--version", "2026.11"); code != 0 || !strings.Contains(out, `format 9, "Apogee" version 2026.11 (uuid `) {
		t.Errorf("create with a title and a version: %d %q %q", code, out, errOut)
	}
	if code, _, errOut := runCmd(t, "create", "--version", " ", filepath.Join(t.TempDir(), "x.sqlite")); code != 1 || !strings.Contains(errOut, "--version must not be empty") {
		t.Errorf("create with an empty version: %d %q", code, errOut)
	}

	empty := filepath.Join(t.TempDir(), "empty-file.sqlite")
	os.WriteFile(empty, nil, 0o644)
	for _, p := range []string{path, empty, t.TempDir()} {
		if code, _, errOut := runCmd(t, "create", p); code != 1 || !strings.Contains(errOut, "already exists") {
			t.Errorf("create over %s: %d %q", p, code, errOut)
		}
	}
}

func TestCheckSmallLibrary(t *testing.T) {
	path := newLibrary(t, smallLibrary...)
	before, _ := os.ReadFile(path)
	code, out, errOut := runCmd(t, "check", path)
	if code != 0 || !strings.Contains(out, "with 1 course, 2 chapters and 4 lessons") {
		t.Errorf("check: %d %q %q", code, out, errOut)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("check changed the file")
	}
}

// Each rule of the format fails on a library broken for it, with one line per problem.
func TestCheckFindsEachProblem(t *testing.T) {
	for _, tc := range []struct {
		name    string
		breakIt []string // statements that break the small library
		want    []string // the lines check prints
	}{
		{
			name:    "another kind of file",
			breakIt: []string{`PRAGMA application_id = 7`},
			want:    []string{"application_id is 7, not 1196641619: the file is not a Gosmos library"},
		},
		{
			name:    "a library of format 6",
			breakIt: []string{`PRAGMA application_id = 1196447556`, `PRAGMA user_version = 6`},
			want:    []string{"application_id is 1196447556 and user_version 6: the library is of format 6, from before format 9; library migrate brings it up to date"},
		},
		{
			name:    "an older format",
			breakIt: []string{`PRAGMA user_version = 8`},
			want:    []string{"user_version is 8, not 9: the library is of an older format; library migrate brings it up to date"},
		},
		{
			name:    "a newer format",
			breakIt: []string{`PRAGMA user_version = 10`},
			want:    []string{"user_version is 10, not 9: the library is of a newer format than this program knows"},
		},
		{
			name:    "a hint for a lesson that is not there",
			breakIt: []string{`PRAGMA foreign_keys = OFF`, `INSERT INTO lesson_hint VALUES ('l-gone', 0, 'Look closer.')`},
			want:    []string{"foreign key: lesson_hint row 6 points at a missing lesson row"},
		},
		{
			name:    "no library_info row",
			breakIt: []string{`DELETE FROM library_info`},
			want:    []string{"library_info holds 0 rows, not exactly 1"},
		},
		{
			name:    "a lesson without its per-type row",
			breakIt: []string{`DELETE FROM quiz_lesson`},
			want:    []string{"lesson l-quiz (type_choice) has 0 quiz_lesson rows, want 1"},
		},
		{
			name:    "a lesson with another type's row",
			breakIt: []string{`PRAGMA foreign_keys = OFF`, `INSERT INTO quiz_lesson VALUES ('l-table', 'type_choice', 0)`},
			want: []string{"foreign key: quiz_lesson row 2 points at a missing lesson row",
				"lesson l-table (type_code_tests) has 1 quiz_lesson rows, want 0"},
		},
		{
			name:    "a lesson of a type format 8 removed",
			breakIt: []string{`PRAGMA foreign_keys = OFF`, `PRAGMA ignore_check_constraints = ON`, `UPDATE lesson SET type = 'type_interview_chat' WHERE uuid = 'l-quiz'`},
			want: []string{"foreign key: quiz_lesson row 1 points at a missing lesson row",
				"lesson l-quiz (type_interview_chat) has 1 quiz_lesson rows, want 0",
				"lesson l-quiz (type_interview_chat) is of a type format 9 does not have"},
		},
		{
			name:    "a quiz question without its right answer",
			breakIt: []string{`UPDATE quiz_answer SET is_correct = 0`},
			want:    []string{"quiz question 0 of lesson l-quiz has 0 right answers, want 1"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := newLibrary(t, append(smallLibrary, tc.breakIt...)...)
			before, _ := os.ReadFile(path)
			code, out, errOut := runCmd(t, "check", path)
			if code != 1 || errOut != "" {
				t.Fatalf("check: exit %d, stderr %q, want 1 and nothing", code, errOut)
			}
			if got := strings.Split(strings.TrimSuffix(out, "\n"), "\n"); strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
				t.Errorf("check printed\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(tc.want, "\n"))
			}
			if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
				t.Error("check changed the file")
			}
		})
	}
}

func TestCheckDamagedAndForeignFiles(t *testing.T) {
	// A page overwritten as a failing disk might: every line comes from quick_check.
	path := newLibrary(t, smallLibrary...)
	d, _, err := db.Inspect(path)
	if err != nil {
		t.Fatal(err)
	}
	var root, size int64
	d.QueryRow(`SELECT (SELECT rootpage FROM sqlite_schema WHERE name = 'quiz_answer'), (SELECT page_size FROM pragma_page_size)`).Scan(&root, &size)
	d.Close()
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteAt(bytes.Repeat([]byte{0xff}, 16), (root-1)*size)
	f.Close()
	code, out, _ := runCmd(t, "check", path)
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if code != 1 || len(lines) == 0 {
		t.Fatalf("check of a damaged file: %d %q", code, out)
	}
	for _, l := range lines {
		if !strings.HasPrefix(l, "quick_check: ") {
			t.Errorf("line %q does not come from quick_check", l)
		}
	}

	notes := filepath.Join(t.TempDir(), "notes.sqlite")
	os.WriteFile(notes, []byte("Shopping: milk, eggs.\n"), 0o644)
	if code, out, _ := runCmd(t, "check", notes); code != 1 || !strings.HasPrefix(out, "the file is not an SQLite database SQLite can read: ") {
		t.Errorf("check of a text file: %d %q", code, out)
	}
	if code, _, errOut := runCmd(t, "check", filepath.Join(t.TempDir(), "missing.sqlite")); code != 1 || !strings.Contains(errOut, "no database file") {
		t.Errorf("check of a missing file: %d %q", code, errOut)
	}
}

// historySchema returns the schema of an older format, as db/history keeps it.
func historySchema(t *testing.T, v int) string {
	t.Helper()
	schema, err := os.ReadFile(filepath.Join("..", "..", "db", "history", fmt.Sprintf("schema_v%d.sql", v)))
	if err != nil {
		t.Fatal(err)
	}
	return string(schema)
}

// v3Library makes a library of format 3 from the schema kept in db/history: a course with a
// code lesson, the learner's progress in it, and the learner's project folder for the course,
// which only command-line lessons used.
func v3Library(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "old.sqlite")
	exec(t, path,
		historySchema(t, 3),
		`INSERT INTO course VALUES ('c-ship', 'learn-shipping', 'Learn Shipping', 'Course', 'Ship it.', 'An invented course about shipping code.', 0)`,
		`INSERT INTO chapter VALUES ('ch-build', 'c-ship', 0, '1-build', 'Build', '')`,
		`INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod, code_last_mod, synced_at)
			VALUES ('l-build', 'ch-build', 0, '1-build', 'Build your app', 'type_code', 1, 30, 'Build it.', '2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', NULL, '2026-01-03T00:00:00Z')`,
		`INSERT INTO code_lesson VALUES ('l-build', 'type_code', 'go', 'built' || char(10))`,
		`INSERT INTO starter_file VALUES ('l-build', 'main.go', 0, 'package main', 0, 0)`,
		`INSERT INTO progress VALUES ('l-build', NULL, '2026-10-31T12:00:00Z')`,
		`INSERT INTO project VALUES ('c-ship', '/Users/me/shipping', 'https://my-app.example.com')`,
	)
	return path
}

// v7Library makes a library of format 7 from the schema kept in db/history: the small library,
// and the statements given.
func v7Library(t *testing.T, stmts ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "v7.sqlite")
	exec(t, path, append([]string{historySchema(t, 7),
		`INSERT INTO library_info VALUES (1, '0b7f3e52-9a41-4c6d-8f20-6e5d1c9b7a34', 'Learn Testing', '2')`}, append(smallLibrary, stmts...)...)...)
	return path
}

// Format 7 goes to 9 in place: the library keeps its uuid, its title and every lesson.
func TestMigrateFrom7(t *testing.T) {
	path := v7Library(t)
	code, out, errOut := runCmd(t, "migrate", path)
	want := `^migrated ` + regexp.QuoteMeta(path) + ` from format 7 to 9: "Learn Testing" version 2 \(uuid 0b7f3e52-9a41-4c6d-8f20-6e5d1c9b7a34\)
ok: ` + regexp.QuoteMeta(path) + ` is "Learn Testing" version 2 \(uuid 0b7f3e52-9a41-4c6d-8f20-6e5d1c9b7a34\), a library of format 9 with 1 course, 2 chapters and 4 lessons
$`
	if code != 0 || !regexp.MustCompile(want).MatchString(out) {
		t.Fatalf("migrate: %d\n%s\n%s", code, out, errOut)
	}
}

// Format 8 goes to 9 in place: the library keeps its uuid, its title and every lesson, and has
// the tables of a course's story, empty. Migrating it again changes nothing.
func TestMigrateFrom8(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v8.sqlite")
	exec(t, path, append([]string{historySchema(t, 8),
		`INSERT INTO library_info VALUES (1, '5d2a9c47-1e8b-4f30-9a6d-3b7c0e2f8a15', 'Learn Testing', '3')`}, smallLibrary...)...)
	code, out, errOut := runCmd(t, "migrate", path)
	want := `^migrated ` + regexp.QuoteMeta(path) + ` from format 8 to 9: "Learn Testing" version 3 \(uuid 5d2a9c47-1e8b-4f30-9a6d-3b7c0e2f8a15\)
ok: ` + regexp.QuoteMeta(path) + ` is "Learn Testing" version 3 \(uuid 5d2a9c47-1e8b-4f30-9a6d-3b7c0e2f8a15\), a library of format 9 with 1 course, 2 chapters and 4 lessons
$`
	if code != 0 || !regexp.MustCompile(want).MatchString(out) {
		t.Fatalf("migrate: %d\n%s\n%s", code, out, errOut)
	}
	d, _, err := db.Inspect(path)
	if err != nil {
		t.Fatal(err)
	}
	var tables, rows int
	d.QueryRow(`SELECT (SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name IN ('story_page', 'story_card', 'chapter_story', 'story_term')),
		(SELECT count(*) FROM story_page) + (SELECT count(*) FROM story_card) + (SELECT count(*) FROM chapter_story) + (SELECT count(*) FROM story_term)`).Scan(&tables, &rows)
	d.Close()
	if tables != 4 || rows != 0 {
		t.Errorf("after migrating: %d of the story's 4 tables, holding %d rows; want all 4, empty", tables, rows)
	}
	before, _ := os.ReadFile(path)
	if code, out, _ := runCmd(t, "migrate", path); code != 0 || !strings.HasPrefix(out, path+` is already of format 9: "Learn Testing" version 3`) {
		t.Errorf("migrating again: %d %q", code, out)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("migrating an up-to-date library changed it")
	}
}

func TestMigrate(t *testing.T) {
	path := v3Library(t)
	if code, out, _ := runCmd(t, "check", path); code != 1 || !strings.HasPrefix(out, "application_id is 1196447556 and user_version 3: the library is of format 3") {
		t.Fatalf("check before migrating: %d %q", code, out)
	}
	unchanged, _ := os.ReadFile(path)
	if code, _, errOut := runCmd(t, "migrate", path); code != 1 || !strings.Contains(errOut, "holds the learner's data (1 progress rows); give a learner database to move it into (--learner <file>): nothing changed") {
		t.Fatalf("migrate without a learner database: %d %q", code, errOut)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(unchanged, after) {
		t.Fatal("the refused migration changed the library")
	}
	learner := filepath.Join(t.TempDir(), "Gosmos", "learner.sqlite")
	code, out, errOut := runCmd(t, "migrate", path, "--learner", learner, "--title", "Shipping")
	want := `^migrated ` + regexp.QuoteMeta(path) + ` from format 3 to 9: "Shipping" version 1 \(uuid ([0-9a-f-]{36})\)
moved its learner data into ` + regexp.QuoteMeta(learner) + `, a new learner database: progress 1, attempt 0, note 0, code_file 0, setting 0
dropped the learner data only command-line lessons used, which Gosmos no longer has: project 1
ok: `
	m := regexp.MustCompile(want).FindStringSubmatch(out)
	if code != 0 || m == nil {
		t.Fatalf("migrate: %d\n%s\n%s", code, out, errOut)
	}
	l, err := db.OpenLearner(learner)
	if err != nil {
		t.Fatal(err)
	}
	var opened string
	l.QueryRow(`SELECT last_opened_at FROM progress WHERE library_uuid = ? AND lesson_uuid = 'l-build'`, m[1]).Scan(&opened)
	l.Close()
	if opened != "2026-10-31T12:00:00Z" {
		t.Errorf("the progress became %q in the learner database", opened)
	}

	before, _ := os.ReadFile(path)
	if code, out, _ := runCmd(t, "migrate", path, "--learner", learner); code != 0 || !strings.HasPrefix(out, path+" is already of format 9: \"Shipping\" version 1 (uuid "+m[1]+")") {
		t.Errorf("migrating again: %d %q", code, out)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("migrating an up-to-date library changed it")
	}
}

func TestMigrateRefuses(t *testing.T) {
	newer := newLibrary(t, `PRAGMA user_version = 10`)
	other := filepath.Join(t.TempDir(), "other.sqlite")
	exec(t, other, `CREATE TABLE recipe (name TEXT)`)
	interview := v7Library(t, `INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
		VALUES ('l-talk', 'ch-more', 1, 'l-talk', 'Talk it through', 'type_interview_chat', 1, 3, 'Explain table tests.', '2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z')`,
		`INSERT INTO interview_lesson VALUES ('l-talk', 'type_interview_chat', 'Explain table-driven tests.', 'Mentions t.Run.')`)
	for _, tc := range []struct{ path, want string }{
		{newer, "newer than this program's 9"},
		{interview, "holds lessons of types format 8 removed (type_interview_chat 1): Gosmos is self-contained, with code, tests and quiz lessons only, which it checks on the Mac alone; nothing changed"},
		{other, "not a Gosmos library (application_id 0); nothing changed"},
		{filepath.Join(t.TempDir(), "missing.sqlite"), "no database file"},
	} {
		before, _ := os.ReadFile(tc.path)
		if code, _, errOut := runCmd(t, "migrate", tc.path); code != 1 || !strings.Contains(errOut, tc.want) {
			t.Errorf("migrate %s: %d %q, want it to say %q", tc.path, code, errOut, tc.want)
		}
		if after, _ := os.ReadFile(tc.path); !bytes.Equal(before, after) {
			t.Errorf("a refused migration changed %s", tc.path)
		}
	}
}

func TestUsage(t *testing.T) {
	for _, args := range [][]string{nil, {"check"}, {"frobnicate", "x.sqlite"}, {"check", "a", "b"}} {
		if code, _, errOut := runCmd(t, args...); code != 2 || !strings.Contains(errOut, "usage:") {
			t.Errorf("run %v: %d %q", args, code, errOut)
		}
	}
}

// crashEnv names, for the child process crash starts, the library to interrupt a write on.
const crashEnv = "GOSMOS_INTERRUPTED_WRITE"

// TestHelperInterruptedWrite is not a test of its own: crash runs it in a child process. It
// starts a write on the library its environment names, too large for SQLite to keep in memory,
// so the write spills to the file before any commit, says so, and waits to be killed.
func TestHelperInterruptedWrite(t *testing.T) {
	path := os.Getenv(crashEnv)
	if path == "" {
		t.Skip("crash runs this in a child process")
	}
	d, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	d.SetMaxOpenConns(1)
	if _, err := d.Exec(`PRAGMA cache_size = 10`); err != nil {
		t.Fatal(err)
	}
	tx, err := d.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`UPDATE course SET title = 'Changed by a write that never finished'`); err != nil {
		t.Fatal(err)
	}
	for i := range 400 {
		if _, err := tx.Exec(`INSERT INTO lesson_hint VALUES ('l-quiz', ?, ?)`, i, fmt.Sprintf("Unfinished hint %03d: %s", i, strings.Repeat("x", 4000))); err != nil {
			t.Fatal(err)
		}
	}
	fmt.Println("spilled")
	time.Sleep(time.Minute) // killed long before this ends
}

// crash interrupts a write on the library at path as a crash or Ctrl-C would, leaving its
// hot journal.
func crash(t *testing.T, path string) {
	t.Helper()
	cmd := osexec.Command(os.Args[0], "-test.run=^TestHelperInterruptedWrite$")
	cmd.Env = append(os.Environ(), crashEnv+"="+path)
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	spilled := make(chan bool, 1)
	go func() {
		for sc := bufio.NewScanner(out); sc.Scan(); {
			if sc.Text() == "spilled" {
				spilled <- true
				return
			}
		}
		spilled <- false
	}()
	select {
	case ok := <-spilled:
		if !ok {
			t.Fatal("the writer ended before its write spilled to the file")
		}
	case <-time.After(time.Minute):
		t.Fatal("the writer did not spill in time")
	}
	cmd.Process.Kill()
	cmd.Wait()
	if info, err := os.Stat(path + "-journal"); err != nil || info.Size() == 0 {
		t.Fatalf("no hot journal after the crash: %v", err)
	}
}

func TestCheckRecoversAnInterruptedWrite(t *testing.T) {
	path := newLibrary(t, smallLibrary...)
	_, ok, _ := runCmd(t, "check", path)
	crash(t, path)

	code, out, errOut := runCmd(t, "check", path)
	if want := db.RecoveredMessage(path) + "\n" + ok; code != 0 || out != want {
		t.Fatalf("check after a crash: exit %d, stderr %q, printed\n%s\nwant\n%s", code, errOut, out, want)
	}
	if _, err := os.Stat(path + "-journal"); err == nil {
		t.Error("the journal is still there")
	}
	d, _, err := db.Inspect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var title string
	var unfinished int
	d.QueryRow(`SELECT (SELECT title FROM course), (SELECT count(*) FROM lesson_hint WHERE text LIKE 'Unfinished hint%')`).Scan(&title, &unfinished)
	if title != "Learn Testing" || unfinished != 0 {
		t.Errorf("after recovering: course title %q and %d unfinished hints; want the state before the write", title, unfinished)
	}
	if code, again, _ := runCmd(t, "check", path); code != 0 || again != ok {
		t.Errorf("a second check: %d %q, want %q", code, again, ok)
	}
}

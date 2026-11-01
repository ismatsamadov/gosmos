package db

// Migrations take a library of an older format to this program's, one format at a time, each
// step in one transaction. A table a step changes gets the definition its new format has for
// it: the step creates it from that format's schema (history/schema_v<N>.sql, or schema.sql
// for this program's format) and copies the rows over.
//
// Format 7 moved the learner's data (progress, attempts, notes, code, projects, server
// addresses and settings) out of the library into a learner database (learner.sql). The step
// from 6 to 7 therefore writes two files: the learner database first, and the library only
// once that is committed.

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed history/*.sql
var history embed.FS

// formatSchema returns the schema of format v: schema.sql for this program's format, the kept
// copy for an older one.
func formatSchema(v int) (string, error) {
	if v == Version {
		return Schema, nil
	}
	b, err := history.ReadFile(fmt.Sprintf("history/schema_v%d.sql", v))
	if err != nil {
		return "", fmt.Errorf("no schema of format %d is kept", v)
	}
	return string(b), nil
}

// oldest is the oldest format Migrate takes.
const oldest = 3

// steps[v] takes a library of format v to v+1 inside the file, in the transaction it is given,
// with foreign keys off. The step from 6 to 7 is Migrate's own: it needs the learner database.
var steps = map[int]func(ctx context.Context, tx *sql.Tx) error{
	3: migrate3to4,
	4: migrate4to5,
	5: migrate5to6,
}

// learnerTables are the tables of format 6 and older that hold the learner's data, in the order
// Migrate moves them.
var learnerTables = []string{"progress", "attempt", "note", "code_file", "project", "server_address", "setting"}

// migrate3to4 makes format 4:
//   - cli_lesson.base_url_default holds only real addresses. The source's 'override', which
//     means the lesson needs the address of your own server, becomes asks_for_base_url = 1.
//   - project.dir may be NULL, so a course can have an address without a folder.
//   - server_address is new.
func migrate3to4(ctx context.Context, tx *sql.Tx) error {
	if err := rebuild(ctx, tx, 4, "cli_lesson", `
		INSERT INTO new_cli_lesson (lesson_uuid, lesson_type, base_url_default, asks_for_base_url, use_agent_prompt, contains_complete_dir)
		SELECT lesson_uuid, lesson_type, nullif(base_url_default, 'override'), base_url_default IS 'override',
		       use_agent_prompt, contains_complete_dir
		FROM cli_lesson`); err != nil {
		return err
	}
	if err := rebuild(ctx, tx, 4, "project", `
		INSERT INTO new_project (course_uuid, dir, base_url) SELECT course_uuid, dir, base_url FROM project`); err != nil {
		return err
	}
	return createTable(ctx, tx, 4, "server_address")
}

// migrate4to5 makes format 5: interview_lesson.acceptance_criteria is NULL where the lesson
// gives none, instead of an empty or blank string. Every other value is copied exactly.
func migrate4to5(ctx context.Context, tx *sql.Tx) error {
	return rebuild(ctx, tx, 5, "interview_lesson", `
		INSERT INTO new_interview_lesson (lesson_uuid, lesson_type, student_prompt, acceptance_criteria)
		SELECT lesson_uuid, lesson_type, student_prompt,
		       CASE WHEN trim(acceptance_criteria, ' ' || char(9, 10, 11, 12, 13)) = '' THEN NULL ELSE acceptance_criteria END
		FROM interview_lesson`)
}

// migrate5to6 makes format 6: solution_file, the reference solutions producers verify their
// lessons with, is new.
func migrate5to6(ctx context.Context, tx *sql.Tx) error {
	return createTable(ctx, tx, 6, "solution_file")
}

// MigrateOptions says how Migrate brings a library of format 6 or older to format 7.
type MigrateOptions struct {
	// Title is the library's title (library_info.title). "" takes the title of its first
	// course, or else the file's name.
	Title string
	// Learner is the learner database to move the library's learner data into; Migrate makes
	// it from learner.sql if it does not exist. It is required when there is any such data.
	Learner string
	// Now is when the library joins the learner database (library.added_at); zero means now.
	Now time.Time
}

// TableRows is how many rows of a table something holds.
type TableRows struct {
	Table string
	Rows  int
}

// Migration says what Migrate did.
type Migration struct {
	From, To  int
	Recovered bool        // an interrupted write was rolled back first (RecoveredMessage)
	Info      Info        // the library's library_info row, once it is of format 7
	Learner   string      // the learner database its learner data went to, if one was given
	Created   bool        // Migrate made that learner database from learner.sql
	Moved     []TableRows // the learner's rows of each table, which the learner database now holds
	Already   bool        // an earlier run had copied them, so none was copied again
}

// Migrate brings the library at path up to this program's format, one format at a time. Each
// step runs in one transaction, with foreign keys checked before it commits.
//
// The step from 6 to 7 gives the library its library_info row (a new uuid, the title from
// opt.Title or else the first course's or the file's name, and version 1), puts lesson
// difficulties on format 7's scale of 1 to 10 (old score + 6, divided by 7), and moves the
// learner's data into the learner database opt.Learner, under a library row for the file
// (location: its absolute path). That transaction commits first. Only then is the library
// rebuilt without the learner tables, in one transaction that checks every other table kept
// its rows, and vacuumed, so no copy of the learner's data stays in it. Run again after an
// interruption, Migrate finds the file's library row in the learner database, takes its uuid
// and title, and copies nothing a second time as long as the row counts match.
//
// A library that holds learner data is refused without opt.Learner, and so is one whose
// learner database is not one of this program's, before anything changes. A file of this
// format is left alone; a newer one, or one that is not a library, is refused.
func Migrate(ctx context.Context, path string, opt MigrateOptions) (Migration, error) {
	m := Migration{To: Version}
	if info, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return m, fmt.Errorf("opening %s: %w", path, ErrNotFound)
	} else if err != nil {
		return m, fmt.Errorf("opening %s: %w", path, err)
	} else if info.IsDir() {
		return m, fmt.Errorf("opening %s: it is a folder", path)
	}
	var err error
	if m.Recovered, err = Recover(path); err != nil {
		return m, err
	}
	d, err := connect(path)
	if err != nil {
		return m, err
	}
	defer d.Close()
	var app int
	if err := d.QueryRowContext(ctx, `SELECT (SELECT application_id FROM pragma_application_id), (SELECT user_version FROM pragma_user_version)`).Scan(&app, &m.From); err != nil {
		return m, fmt.Errorf("%s: %w", path, err)
	}
	switch {
	case app == ApplicationID && m.From == Version:
		m.Info, err = ReadInfo(ctx, d)
		return m, err
	case app == LegacyApplicationID && m.From >= oldest && m.From <= 6:
	case app == LegacyApplicationID && m.From < oldest:
		return m, fmt.Errorf("%s is of format %d; library migrate takes format %d and later only", path, m.From, oldest)
	default:
		if err := formatError(app, m.From); err != nil {
			return m, fmt.Errorf("%s: %w; nothing changed", path, err)
		}
		return m, fmt.Errorf("%s: no migration from format %d", path, m.From)
	}

	// Everything that can refuse the migration does so before it changes anything.
	held, err := learnerRows(ctx, d)
	if err != nil {
		return m, err
	}
	if total(held) > 0 && opt.Learner == "" {
		return m, fmt.Errorf("%s holds the learner's data (%s); give a learner database to move it into (--learner <file>): nothing changed", path, describe(held))
	}
	var l *sql.DB
	if opt.Learner != "" {
		if same, err := sameFile(path, opt.Learner); err != nil || same {
			return m, errors.Join(err, fmt.Errorf("the learner database cannot be the library itself: nothing changed"))
		}
		if _, err := os.Stat(opt.Learner); err == nil {
			if l, err = OpenLearner(opt.Learner); err != nil {
				return m, fmt.Errorf("%w: nothing changed", err)
			}
			defer l.Close()
		}
	}

	for v := m.From; v < 6; v++ {
		if err := migrate(ctx, d, v, steps[v]); err != nil {
			return m, err
		}
	}
	info := Info{UUID: NewUUID(), Version: "1"}
	if info.Title, err = pickTitle(ctx, d, path, opt.Title); err != nil {
		return m, err
	}
	if opt.Learner != "" {
		if l == nil {
			if l, err = CreateLearner(opt.Learner); err != nil {
				return m, err
			}
			defer l.Close()
			m.Created = true
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return m, err
		}
		m.Learner = opt.Learner
		if m.Moved, m.Already, err = moveLearnerData(ctx, l, abs, &info, opt.Now); err != nil {
			return m, err
		}
	}
	if err := migrate(ctx, d, 6, func(ctx context.Context, tx *sql.Tx) error { return migrate6to7(ctx, tx, info) }); err != nil {
		return m, err
	}
	// The learner tables' pages are free now but still hold their rows: VACUUM writes the
	// file again without them. Lessons keep their ids (lesson.id is an INTEGER PRIMARY KEY).
	if _, err := d.ExecContext(ctx, `VACUUM`); err != nil {
		return m, fmt.Errorf("vacuuming %s: %w", path, err)
	}
	m.Info = info
	return m, nil
}

// migrate6to7 makes format 7 inside the library, once its learner data is in the learner
// database: the learner tables go; lesson loses code_last_mod and synced_at, and its difficulty
// goes from the old scores, about 5 to 70, to 1 to 10; code_lesson gains solution_note and
// quiz_answer explanation; library_info (holding info) and lesson_hint are new; and the file
// gets format 7's mark. Every table it keeps must hold as many rows as before, and the search
// index must still match the lessons.
func migrate6to7(ctx context.Context, tx *sql.Tx, info Info) error {
	before, err := tableRows(ctx, tx)
	if err != nil {
		return err
	}
	for _, t := range learnerTables {
		if _, err := tx.ExecContext(ctx, `DROP TABLE `+t); err != nil {
			return fmt.Errorf("dropping %s: %w", t, err)
		}
		delete(before, t)
	}
	// Each lesson keeps its id, the search index's key.
	if err := rebuild(ctx, tx, 7, "lesson", `
		INSERT INTO new_lesson (id, uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
		SELECT id, uuid, chapter_uuid, position, slug, title, type, required, min(10, max(1, (difficulty + 6) / 7)), readme, created_at, last_mod
		FROM lesson`); err != nil {
		return err
	}
	if err := rebuild(ctx, tx, 7, "code_lesson", `
		INSERT INTO new_code_lesson (lesson_uuid, lesson_type, prog_lang, expected_output)
		SELECT lesson_uuid, lesson_type, prog_lang, expected_output FROM code_lesson`); err != nil {
		return err
	}
	if err := rebuild(ctx, tx, 7, "quiz_answer", `
		INSERT INTO new_quiz_answer (lesson_uuid, question_position, position, answer, is_correct)
		SELECT lesson_uuid, question_position, position, answer, is_correct FROM quiz_answer`); err != nil {
		return err
	}
	for _, t := range []string{"library_info", "lesson_hint"} {
		if err := createTable(ctx, tx, 7, t); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO library_info (singleton, uuid, title, version) VALUES (1, ?, ?, ?)`, info.UUID, info.Title, info.Version); err != nil {
		return fmt.Errorf("library_info: %w", err)
	}
	after, err := tableRows(ctx, tx)
	if err != nil {
		return err
	}
	for t, n := range before {
		if after[t] != n {
			return fmt.Errorf("%s holds %d rows after the rebuild, not %d", t, after[t], n)
		}
	}
	if after["library_info"] != 1 || after["lesson_hint"] != 0 {
		return fmt.Errorf("library_info holds %d rows and lesson_hint %d, not 1 and 0", after["library_info"], after["lesson_hint"])
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO lesson_search(lesson_search, rank) VALUES ('integrity-check', 1)`); err != nil {
		return fmt.Errorf("the search index no longer matches the lessons: %w", err)
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf(`PRAGMA application_id = %d`, ApplicationID)); err != nil {
		return err
	}
	return nil
}

// moveLearnerData copies the learner's data from the library file at abs into the learner
// database l, in one transaction: a library row for the file, then its rows of every learner
// table, keyed by the library's uuid. Attempts get new ids in their old order. A server address
// or setting the learner database already has keeps its value there. Every table's rows are
// counted on both sides before the transaction commits.
//
// A library row already there for the same location is an earlier run's: its uuid and title
// become info's, and as long as the counts match, nothing is copied again.
func moveLearnerData(ctx context.Context, l *sql.DB, abs string, info *Info, now time.Time) ([]TableRows, bool, error) {
	conn, err := l.Conn(ctx)
	if err != nil {
		return nil, false, err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, `ATTACH DATABASE ? AS content`, abs); err != nil {
		return nil, false, fmt.Errorf("reading %s from the learner database: %w", abs, err)
	}
	defer conn.ExecContext(context.Background(), `DETACH DATABASE content`)
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	var added string
	err = tx.QueryRowContext(ctx, `SELECT uuid, title, added_at FROM main.library WHERE location = ?`, abs).Scan(&info.UUID, &info.Title, &added)
	already := err == nil
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	}
	if !already {
		if now.IsZero() {
			now = time.Now()
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO main.library (uuid, title, location, added_at) VALUES (?, ?, ?, ?)`,
			info.UUID, info.Title, abs, now.UTC().Truncate(time.Second).Format(time.RFC3339)); err != nil {
			return nil, false, fmt.Errorf("the library row: %w", err)
		}
		for _, s := range []string{
			`INSERT INTO main.progress (library_uuid, lesson_uuid, done_at, last_opened_at)
				SELECT ?1, lesson_uuid, done_at, last_opened_at FROM content.progress ORDER BY lesson_uuid`,
			`INSERT INTO main.attempt (library_uuid, lesson_uuid, kind, passed, output, duration_ms, at)
				SELECT ?1, lesson_uuid, kind, passed, output, duration_ms, at FROM content.attempt ORDER BY id`,
			`INSERT INTO main.note (library_uuid, lesson_uuid, body, updated_at)
				SELECT ?1, lesson_uuid, body, updated_at FROM content.note ORDER BY lesson_uuid`,
			`INSERT INTO main.code_file (library_uuid, lesson_uuid, name, content, updated_at)
				SELECT ?1, lesson_uuid, name, content, updated_at FROM content.code_file ORDER BY lesson_uuid, name`,
			`INSERT INTO main.project (library_uuid, course_uuid, dir, base_url)
				SELECT ?1, course_uuid, dir, base_url FROM content.project ORDER BY course_uuid`,
			`INSERT INTO main.server_address (default_url, url)
				SELECT default_url, url FROM content.server_address WHERE ?1 IS NOT NULL ORDER BY default_url
				ON CONFLICT (default_url) DO NOTHING`,
			`INSERT INTO main.setting (key, value)
				SELECT key, value FROM content.setting WHERE ?1 IS NOT NULL ORDER BY key
				ON CONFLICT (key) DO NOTHING`,
		} {
			if _, err := tx.ExecContext(ctx, s, info.UUID); err != nil {
				return nil, false, fmt.Errorf("copying the learner's data: %w", err)
			}
		}
	}

	moved, mismatch, err := compareLearnerRows(ctx, tx, info.UUID)
	if err != nil {
		return nil, false, err
	}
	if len(mismatch) > 0 {
		if already {
			return nil, false, fmt.Errorf("the learner database already has this library (uuid %s, added %s), with other rows than the file: %s; nothing changed",
				info.UUID, added, strings.Join(mismatch, ", "))
		}
		return nil, false, fmt.Errorf("copying the learner's data went wrong (%s); nothing changed", strings.Join(mismatch, ", "))
	}
	if already {
		return moved, true, nil
	}
	return moved, false, tx.Commit()
}

// compareLearnerRows counts, for every learner table, the file's rows and the ones the learner
// database holds of them for the library: the rows of the library's uuid, or for server
// addresses and settings, the rows with the same key. Attempts must also come in the same
// order. It returns the file's counts and a line for each difference.
func compareLearnerRows(ctx context.Context, tx *sql.Tx, uuid string) ([]TableRows, []string, error) {
	held := func(t string) string {
		switch t {
		case "server_address":
			return `SELECT count(*) FROM content.server_address c WHERE EXISTS (SELECT 1 FROM main.server_address s WHERE s.default_url = c.default_url)`
		case "setting":
			return `SELECT count(*) FROM content.setting c WHERE EXISTS (SELECT 1 FROM main.setting s WHERE s.key = c.key)`
		}
		return `SELECT count(*) FROM main.` + t + ` WHERE library_uuid = ?1`
	}
	var counts []TableRows
	var mismatch []string
	for _, t := range learnerTables {
		var want, got int
		if err := tx.QueryRowContext(ctx, `SELECT (SELECT count(*) FROM content.`+t+`), (`+held(t)+`)`, uuid).Scan(&want, &got); err != nil {
			return nil, nil, fmt.Errorf("counting %s: %w", t, err)
		}
		counts = append(counts, TableRows{t, want})
		if got != want {
			mismatch = append(mismatch, fmt.Sprintf("%s: %d rows in the file, %d in the learner database", t, want, got))
		}
	}
	var misplaced int
	if err := tx.QueryRowContext(ctx, `WITH
		a AS (SELECT row_number() OVER (ORDER BY id) AS n, lesson_uuid, kind, passed, output, duration_ms, at FROM content.attempt),
		b AS (SELECT row_number() OVER (ORDER BY id) AS n, lesson_uuid, kind, passed, output, duration_ms, at FROM main.attempt WHERE library_uuid = ?1)
		SELECT count(*) FROM a JOIN b USING (n)
		WHERE a.lesson_uuid IS NOT b.lesson_uuid OR a.kind IS NOT b.kind OR a.passed IS NOT b.passed
		   OR a.output IS NOT b.output OR a.duration_ms IS NOT b.duration_ms OR a.at IS NOT b.at`, uuid).Scan(&misplaced); err != nil {
		return nil, nil, fmt.Errorf("comparing attempts: %w", err)
	}
	if misplaced > 0 {
		mismatch = append(mismatch, fmt.Sprintf("attempt: %d rows differ from the file's, in id order", misplaced))
	}
	return counts, mismatch, nil
}

// learnerRows counts the rows of the learner tables a library of format 6 or older has.
func learnerRows(ctx context.Context, d *sql.DB) ([]TableRows, error) {
	var out []TableRows
	for _, t := range learnerTables {
		var exists, n int
		if err := d.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name = ?`, t).Scan(&exists); err != nil {
			return nil, err
		}
		if exists == 1 {
			if err := d.QueryRowContext(ctx, `SELECT count(*) FROM `+t).Scan(&n); err != nil {
				return nil, fmt.Errorf("counting %s: %w", t, err)
			}
		}
		out = append(out, TableRows{t, n})
	}
	return out, nil
}

func total(rows []TableRows) int {
	n := 0
	for _, r := range rows {
		n += r.Rows
	}
	return n
}

// describe lists the tables that hold rows, such as "3 progress and 1 code_file rows".
func describe(rows []TableRows) string {
	var parts []string
	for _, r := range rows {
		if r.Rows > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", r.Rows, r.Table))
		}
	}
	switch len(parts) {
	case 0:
		return "no rows"
	case 1:
		return parts[0] + " rows"
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1] + " rows"
}

// tableRows counts the rows of every table of the file but the full-text index's virtual one.
func tableRows(ctx context.Context, tx *sql.Tx) (map[string]int, error) {
	rows, err := tx.QueryContext(ctx, `SELECT name FROM sqlite_schema WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND sql NOT LIKE 'CREATE VIRTUAL TABLE%'`)
	if err != nil {
		return nil, err
	}
	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			rows.Close()
			return nil, err
		}
		names = append(names, n)
	}
	rows.Close()
	out := map[string]int{}
	for _, t := range names {
		var n int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM "`+t+`"`).Scan(&n); err != nil {
			return nil, fmt.Errorf("counting %s: %w", t, err)
		}
		out[t] = n
	}
	return out, nil
}

// pickTitle returns the title a migrated library gets: the one given, or else its first
// course's, or else the file's name without its extension.
func pickTitle(ctx context.Context, d *sql.DB, path, given string) (string, error) {
	if t := strings.TrimSpace(given); t != "" {
		return t, nil
	}
	var first string
	err := d.QueryRowContext(ctx, `SELECT title FROM course ORDER BY position LIMIT 1`).Scan(&first)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if t := strings.TrimSpace(first); t != "" {
		return t, nil
	}
	base := filepath.Base(path)
	if name := strings.TrimSuffix(base, filepath.Ext(base)); name != "" {
		return name, nil
	}
	return base, nil
}

// sameFile reports whether two paths name the same existing file.
func sameFile(a, b string) (bool, error) {
	fa, err := os.Stat(a)
	if err != nil {
		return false, err
	}
	fb, err := os.Stat(b)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return os.SameFile(fa, fb), nil
}

// migrate runs a step from format from in one transaction, on one connection, with foreign
// keys off while it runs and foreign_key_check before it commits: steps 1, 2, 10, 11 and 12 of
// SQLite's procedure for schema changes ALTER TABLE cannot make
// (https://www.sqlite.org/lang_altertable.html#otheralter). rebuild does the steps between.
func migrate(ctx context.Context, d *sql.DB, from int, step func(context.Context, *sql.Tx) error) (err error) {
	if step == nil {
		return fmt.Errorf("no migration from format %d", from)
	}
	conn, err := d.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Outside a transaction, where the pragma takes effect.
	if _, err := conn.ExecContext(ctx, `PRAGMA foreign_keys = OFF`); err != nil {
		return err
	}
	defer func() {
		if _, on := conn.ExecContext(context.Background(), `PRAGMA foreign_keys = ON`); on != nil {
			err = errors.Join(err, fmt.Errorf("turning foreign keys back on: %w", on))
		}
	}()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := step(ctx, tx); err != nil {
		return fmt.Errorf("migrating from format %d: %w", from, err)
	}
	var broken []string
	rows, err := tx.QueryContext(ctx, `SELECT "table", rowid, parent FROM pragma_foreign_key_check`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var table, parent string
		var rowid sql.NullInt64
		if err := rows.Scan(&table, &rowid, &parent); err != nil {
			rows.Close()
			return err
		}
		broken = append(broken, fmt.Sprintf("%s row %d points at a missing %s row", table, rowid.Int64, parent))
	}
	rows.Close()
	if len(broken) > 0 {
		return fmt.Errorf("migrating from format %d would leave foreign keys pointing nowhere, so nothing changed:\n  %s", from, strings.Join(broken, "\n  "))
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf(`PRAGMA user_version = %d`, from+1)); err != nil {
		return err
	}
	return tx.Commit()
}

// rebuild gives table the definition format has for it, by steps 3 to 8 of SQLite's
// procedure: remember the table's indexes and triggers, create the new table as new_<table>,
// copy the rows with copyRows, drop the old table, rename the new one to its name, and make
// the indexes and triggers again. Foreign keys must be off, so that dropping the table does
// not delete the rows that hang from it; migrate sees to that.
func rebuild(ctx context.Context, tx *sql.Tx, format int, table, copyRows string) error {
	create, err := tableDefinition(format, table)
	if err != nil {
		return err
	}
	head := "CREATE TABLE " + table + " ("
	tmp := "new_" + table
	var keep []string
	rows, err := tx.QueryContext(ctx, `SELECT sql FROM sqlite_schema WHERE tbl_name = ? AND type IN ('index', 'trigger') AND sql IS NOT NULL`, table)
	if err != nil {
		return err
	}
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			rows.Close()
			return err
		}
		keep = append(keep, s)
	}
	rows.Close()
	steps := append([]string{
		"CREATE TABLE " + tmp + " (" + strings.TrimPrefix(create, head),
		copyRows,
		"DROP TABLE " + table,
		"ALTER TABLE " + tmp + " RENAME TO " + table,
	}, keep...)
	for _, s := range steps {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("rebuilding %s: %w", table, err)
		}
	}
	return nil
}

// createTable creates a table as format defines it.
func createTable(ctx context.Context, tx *sql.Tx, format int, table string) error {
	create, err := tableDefinition(format, table)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, create); err != nil {
		return fmt.Errorf("%s: %w", table, err)
	}
	return nil
}

// tableDefinition returns the CREATE TABLE statement the schema of format has for a table.
func tableDefinition(format int, table string) (string, error) {
	schema, err := formatSchema(format)
	if err != nil {
		return "", err
	}
	head := "CREATE TABLE " + table + " ("
	for _, s := range statements(schema) {
		if body := withoutLeadingComments(s); strings.HasPrefix(body, head) {
			return body, nil
		}
	}
	return "", fmt.Errorf("the schema of format %d defines no table %s", format, table)
}

// statements splits SQL text at the semicolons outside comments and quotes. A trigger's body
// is split too, which is fine for finding CREATE TABLE statements.
func statements(text string) []string {
	var out []string
	var b strings.Builder
	var quote byte
	comment := false
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch {
		case comment:
			comment = c != '\n'
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '-' && i+1 < len(text) && text[i+1] == '-':
			comment = true
		case c == '\'' || c == '"':
			quote = c
		case c == ';':
			out = append(out, b.String())
			b.Reset()
			continue
		}
		b.WriteByte(c)
	}
	if strings.TrimSpace(b.String()) != "" {
		out = append(out, b.String())
	}
	return out
}

// withoutLeadingComments drops the blank and comment lines before a statement.
func withoutLeadingComments(s string) string {
	for {
		s = strings.TrimLeft(s, " \t\r\n")
		if !strings.HasPrefix(s, "--") {
			return s
		}
		i := strings.IndexByte(s, '\n')
		if i < 0 {
			return ""
		}
		s = s[i+1:]
	}
}

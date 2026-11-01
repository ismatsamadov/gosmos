// Package db creates, opens, checks and migrates Gosmos course libraries, SQLite files whose
// tables are db/schema.sql, and makes and opens the learner databases of db/learner.sql that
// migrating moves a library's learner data into. Both schemas are embedded here.
package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "modernc.org/sqlite" // the pure-Go SQLite driver, registered as "sqlite"
)

// Schema is db/schema.sql. It makes a new, empty database and is never run against an
// existing one.
//
//go:embed schema.sql
var Schema string

// Version is the format's version, and ApplicationID the mark of a Gosmos course library, as
// schema.sql sets them with PRAGMA user_version and PRAGMA application_id.
var (
	Version       = pragmaValue("schema.sql", Schema, "user_version")
	ApplicationID = pragmaValue("schema.sql", Schema, "application_id")
)

// LegacyApplicationID ('GPWD') marks a library of format 6 or older, from before the app was
// called Gosmos. Only Migrate opens such a file.
const LegacyApplicationID = 1196447556

// ErrNotFound is returned by Open for a file that does not exist.
var ErrNotFound = errors.New("no database file")

func pragmaValue(file, schema, name string) int {
	m := regexp.MustCompile(`(?im)^\s*PRAGMA\s+` + name + `\s*=\s*(\d+)\s*;`).FindStringSubmatch(schema)
	if m == nil {
		panic("db: " + file + " does not set " + name)
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		panic("db: bad " + name + " in " + file + ": " + m[1])
	}
	return n
}

// Info is a library's library_info row: which library the file is.
type Info struct {
	UUID    string // never changes: the learner database refers to the library by it
	Title   string // such as Apogee
	Version string // the content's version, such as 1
}

// NewUUID returns a new random uuid (version 4).
func NewUUID() string {
	var b [16]byte
	rand.Read(b[:]) // never fails: crypto/rand stops the program if it cannot read
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// ReadInfo reads a library's library_info row.
func ReadInfo(ctx context.Context, q Querier) (Info, error) {
	var i Info
	err := q.QueryRowContext(ctx, `SELECT uuid, title, version FROM library_info`).Scan(&i.UUID, &i.Title, &i.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return i, errors.New("the library has no library_info row")
	}
	return i, err
}

// Create makes a new, empty library at path from schema.sql, creating its folder if needed,
// with info as its library_info row; a new uuid if info has none. Nothing may exist at path
// yet.
func Create(path string, info Info) (*sql.DB, error) {
	if _, err := os.Lstat(path); err == nil {
		return nil, fmt.Errorf("creating %s: something already exists there", path)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	if info.UUID == "" {
		info.UUID = NewUUID()
	}
	d, err := connect(path)
	if err != nil {
		return nil, err
	}
	if err := create(d, Schema, func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO library_info (singleton, uuid, title, version) VALUES (1, ?, ?, ?)`, info.UUID, info.Title, info.Version)
		return err
	}); err != nil {
		d.Close()
		os.Remove(path)
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	return d, nil
}

// Open opens the existing library at path to write to it. It must be a Gosmos library of this
// program's format. One of an older format is refused and left as it is: bringing it up to date
// may move learner data into a learner database, which only Migrate (library migrate) does.
func Open(path string) (*sql.DB, error) {
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("opening %s: %w", path, ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("opening %s: it is a folder", path)
	}
	// SQLite would roll back an interrupted write on its own here; doing it first says so.
	if recovered, err := Recover(path); err != nil {
		return nil, err
	} else if recovered {
		log.Print(RecoveredMessage(path))
	}
	d, err := connect(path)
	if err != nil {
		return nil, err
	}
	if err := checkFormat(d); err != nil {
		d.Close()
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	return d, nil
}

// Inspect opens the existing file at path to read it as it is: read-only, without the checks
// Open makes, so that Validate can say what is wrong with any file instead of refusing it. Nothing is written through it. An interrupted write is the one thing set right
// first: a read-only connection cannot read past the journal it left, so Inspect has SQLite
// roll it back (Recover) and reports that it did.
func Inspect(path string) (d *sql.DB, recovered bool, err error) {
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, fmt.Errorf("opening %s: %w", path, ErrNotFound)
	}
	if err != nil {
		return nil, false, fmt.Errorf("opening %s: %w", path, err)
	}
	if info.IsDir() {
		return nil, false, fmt.Errorf("opening %s: it is a folder", path)
	}
	if recovered, err = Recover(path); err != nil {
		return nil, false, err
	}
	if d, err = readOnly(path); err != nil {
		return nil, recovered, err
	}
	if err := probe(d); isReadOnlyRollback(err) { // a journal left between Recover and now
		d.Close()
		again, err := Recover(path)
		if err != nil {
			return nil, recovered, err
		}
		recovered = recovered || again
		if d, err = readOnly(path); err != nil {
			return nil, recovered, err
		}
	}
	return d, recovered, nil // a file SQLite cannot read at all is for Validate to report
}

// readOnly opens a read-only connection pool on the file at path.
func readOnly(path string) (*sql.DB, error) {
	q := url.Values{}
	q.Set("mode", "ro")
	q.Add("_pragma", "busy_timeout(10000)")
	u := url.URL{Scheme: "file", Path: path, RawQuery: q.Encode()}
	if !filepath.IsAbs(path) {
		u = url.URL{Scheme: "file", Opaque: (&url.URL{Path: path}).EscapedPath(), RawQuery: q.Encode()}
	}
	d, err := sql.Open("sqlite", u.String())
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	d.SetMaxOpenConns(1)
	return d, nil
}

// connect opens a connection pool on path in which every connection has foreign keys on and
// uses the rollback journal (journal_mode DELETE, never WAL), so the file is whole, and safe
// to commit, whenever no write is under way. The pool holds a single connection: with a
// rollback journal, a second connection of the same process that is still reading would only
// block the first one's writes.
func connect(path string) (*sql.DB, error) {
	if strings.ContainsRune(path, '?') {
		return nil, fmt.Errorf("database path %q must not contain '?'", path)
	}
	q := url.Values{}
	q.Add("_pragma", "busy_timeout(10000)")
	q.Add("_pragma", "foreign_keys(1)")
	q.Add("_pragma", "journal_mode(DELETE)")
	q.Add("_pragma", "synchronous(FULL)")
	q.Set("_txlock", "immediate") // take the write lock at BEGIN, so a transaction never waits halfway
	d, err := sql.Open("sqlite", path+"?"+q.Encode())
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	d.SetMaxOpenConns(1)
	if err := d.Ping(); err != nil {
		d.Close()
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	return d, nil
}

// checkFormat checks that the file is a Gosmos library of this program's format.
func checkFormat(d *sql.DB) error {
	var mode string
	if err := d.QueryRow(`PRAGMA journal_mode`).Scan(&mode); err != nil {
		return err
	}
	if !strings.EqualFold(mode, "delete") {
		return fmt.Errorf("journal mode is %s, want delete", mode)
	}
	var app, version int
	if err := d.QueryRow(`SELECT (SELECT application_id FROM pragma_application_id), (SELECT user_version FROM pragma_user_version)`).Scan(&app, &version); err != nil {
		return err
	}
	return formatError(app, version)
}

// formatError says why a file with this application_id and user_version is not a library of
// this program's format, or returns nil if it is one.
func formatError(app, version int) error {
	switch {
	case app == ApplicationID && version == Version:
		return nil
	case app == ApplicationID && version > Version:
		return fmt.Errorf("the library is of format %d, newer than this program's %d; update the program", version, Version)
	case app == ApplicationID || app == LegacyApplicationID && version <= 6:
		return fmt.Errorf("the library is of format %d, older than %d; library migrate brings it up to date", version, Version)
	}
	return fmt.Errorf("not a Gosmos library (application_id %d)", app)
}

// create runs a schema on an empty database, and then fill, in one transaction.
func create(d *sql.DB, schema string, fill func(*sql.Tx) error) error {
	var objects int
	if err := d.QueryRow(`SELECT count(*) FROM sqlite_schema`).Scan(&objects); err != nil {
		return err
	}
	if objects > 0 {
		return errors.New("the database is not empty")
	}
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("the schema: %w", err)
	}
	if fill != nil {
		if err := fill(tx); err != nil {
			return err
		}
	}
	return tx.Commit()
}

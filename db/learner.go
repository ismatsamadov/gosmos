package db

// The learner database (learner.sql) holds everything that is the learner's: the libraries
// they use, progress, attempts, notes, code, hints and solutions looked at, and the review
// schedule. The app owns it; the tools here only make one and move a library's learner data
// into it when they migrate the library to format 7.

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// LearnerSchema is db/learner.sql. It makes a new, empty learner database.
//
//go:embed learner.sql
var LearnerSchema string

// LearnerVersion is the learner database's version, and LearnerApplicationID its mark ('GSML'),
// as learner.sql sets them.
var (
	LearnerVersion       = pragmaValue("learner.sql", LearnerSchema, "user_version")
	LearnerApplicationID = pragmaValue("learner.sql", LearnerSchema, "application_id")
)

// CreateLearner makes a new, empty learner database at path from learner.sql, creating its
// folder if needed. Nothing may exist at path yet.
func CreateLearner(path string) (*sql.DB, error) {
	if _, err := os.Lstat(path); err == nil {
		return nil, fmt.Errorf("creating %s: something already exists there", path)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	d, err := connectLearner(path)
	if err != nil {
		return nil, err
	}
	if err := create(d, LearnerSchema, nil); err != nil {
		d.Close()
		os.Remove(path)
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	return d, nil
}

// OpenLearner opens the existing learner database at path. It must be a Gosmos learner
// database of this program's version; any other file is refused and left as it is.
func OpenLearner(path string) (*sql.DB, error) {
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
	d, err := connectLearner(path) // a read-write connection rolls back a hot journal itself
	if err != nil {
		return nil, err
	}
	var app, version int
	err = d.QueryRow(`SELECT (SELECT application_id FROM pragma_application_id), (SELECT user_version FROM pragma_user_version)`).Scan(&app, &version)
	switch {
	case err != nil:
		err = fmt.Errorf("opening %s: %w", path, err)
	case app != LearnerApplicationID:
		err = fmt.Errorf("%s is not a Gosmos learner database (application_id %d, not %d)", path, app, LearnerApplicationID)
	case version != LearnerVersion:
		err = fmt.Errorf("%s is a learner database of version %d; this program writes version %d only", path, version, LearnerVersion)
	}
	if err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

// connectLearner opens a connection pool on a learner database with foreign keys on. Unlike a
// library, it keeps whatever journal mode the file has: the app may use WAL for it.
func connectLearner(path string) (*sql.DB, error) {
	if strings.ContainsRune(path, '?') {
		return nil, fmt.Errorf("database path %q must not contain '?'", path)
	}
	q := url.Values{}
	q.Add("_pragma", "busy_timeout(10000)")
	q.Add("_pragma", "foreign_keys(1)")
	q.Add("_pragma", "synchronous(FULL)")
	q.Set("_txlock", "immediate")
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

package db

// The Go reference (reference.sql) is the documentation the app ships: the standard library
// and five articles from go.dev, in one SQLite file that cmd/reference makes and the app only
// reads.

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ReferenceSchema is db/reference.sql. It makes a new, empty reference database.
//
//go:embed reference.sql
var ReferenceSchema string

// ReferenceVersion is the reference database's version, and ReferenceApplicationID its mark
// ('GSMR'), as reference.sql sets them.
var (
	ReferenceVersion       = pragmaValue("reference.sql", ReferenceSchema, "user_version")
	ReferenceApplicationID = pragmaValue("reference.sql", ReferenceSchema, "application_id")
)

// CreateReference makes a new, empty reference database at path from reference.sql, creating
// its folder if needed. Nothing may exist at path yet.
func CreateReference(path string) (*sql.DB, error) {
	if _, err := os.Lstat(path); err == nil {
		return nil, fmt.Errorf("creating %s: something already exists there", path)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	d, err := connect(path)
	if err != nil {
		return nil, err
	}
	if err := create(d, ReferenceSchema, nil); err != nil {
		d.Close()
		os.Remove(path)
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	return d, nil
}

// OpenReference opens the existing reference database at path read-only. It must be a Gosmos
// reference of this program's version; any other file is refused.
func OpenReference(path string) (*sql.DB, error) {
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
	d, err := readOnly(path)
	if err != nil {
		return nil, err
	}
	var app, version int
	err = d.QueryRow(`SELECT (SELECT application_id FROM pragma_application_id), (SELECT user_version FROM pragma_user_version)`).Scan(&app, &version)
	switch {
	case err != nil:
		err = fmt.Errorf("opening %s: %w", path, err)
	case app != ReferenceApplicationID:
		err = fmt.Errorf("%s is not a Gosmos Go reference (application_id %d, not %d)", path, app, ReferenceApplicationID)
	case version != ReferenceVersion:
		err = fmt.Errorf("%s is a Go reference of version %d; this program reads version %d only", path, version, ReferenceVersion)
	}
	if err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

// reLink finds the links of the reference's Markdown: [text](godoc:…) and
// [text](godoc-article:…). Any other kind of link is a problem of its own.
var reLink = regexp.MustCompile(`\]\(([a-z][a-z0-9+.-]*:[^)\s]*)\)`)

// ValidateReference runs the checks docs/reference-format.md lists, and returns one line for
// each problem:
//
//  1. application_id is 1196641618 ('GSMR') and user_version is this program's version;
//  2. quick_check is ok and foreign_key_check finds nothing;
//  3. source holds exactly one row;
//  4. every method's receiver is a type of its package, and every declaration has a symbol;
//  5. search_entry has one row for each package, symbol and article section, and the search
//     index holds them all;
//  6. every link is a godoc: or godoc-article: link to something in the file.
//
// Like Validate, the checks after the first run only once it passes, and only on a file
// quick_check finds whole. ValidateReference writes nothing. An error means a check could not
// run at all.
func ValidateReference(ctx context.Context, d Querier) (Problems, error) {
	var p Problems
	add := func(format string, args ...any) { p = append(p, fmt.Sprintf(format, args...)) }

	var app, version int
	if err := d.QueryRowContext(ctx, `SELECT (SELECT application_id FROM pragma_application_id), (SELECT user_version FROM pragma_user_version)`).Scan(&app, &version); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		add("the file is not an SQLite database SQLite can read: %v", err)
		return p, nil
	}
	switch {
	case app != ReferenceApplicationID:
		add("application_id is %d, not %d: the file is not a Gosmos Go reference", app, ReferenceApplicationID)
	case version != ReferenceVersion:
		add("user_version is %d, not %d: the reference is of a version this program does not read", version, ReferenceVersion)
	}
	if len(p) > 0 {
		return p, nil
	}

	if err := eachRow(ctx, d, `SELECT quick_check FROM pragma_quick_check`, func(r *sql.Rows) error {
		var msg string
		if err := r.Scan(&msg); err != nil {
			return err
		}
		for line := range strings.SplitSeq(msg, "\n") {
			if line = strings.TrimSpace(line); line != "" && line != "ok" && !strings.HasPrefix(line, "*** in database") {
				add("quick_check: %s", line)
			}
		}
		return nil
	}); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		add("quick_check: %v", err)
	}
	if len(p) > 0 {
		return p, nil
	}

	if err := eachRow(ctx, d, `SELECT "table", rowid, parent FROM pragma_foreign_key_check ORDER BY 1, 2`, func(r *sql.Rows) error {
		var table, parent string
		var rowid sql.NullInt64
		if err := r.Scan(&table, &rowid, &parent); err != nil {
			return err
		}
		add("foreign key: %s row %d points at a missing %s row", table, rowid.Int64, parent)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("foreign_key_check: %w", err)
	}

	var sources int
	if err := d.QueryRowContext(ctx, `SELECT count(*) FROM source`).Scan(&sources); err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	if sources != 1 {
		add("source holds %d rows, not exactly 1", sources)
	}

	if err := eachRow(ctx, d, `SELECT p.import_path, s.receiver, s.name FROM symbol s JOIN package p ON p.id = s.package_id
		WHERE s.kind = 'method' AND NOT EXISTS (SELECT 1 FROM symbol t
			WHERE t.package_id = s.package_id AND t.kind = 'type' AND t.name = s.receiver AND t.receiver IS NULL)
		ORDER BY 1, 2, 3`, func(r *sql.Rows) error {
		var path, recv, name string
		if err := r.Scan(&path, &recv, &name); err != nil {
			return err
		}
		add("method %s.%s.%s: its package has no type %s", path, recv, name, recv)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("receivers: %w", err)
	}
	if err := eachRow(ctx, d, `SELECT d.id, p.import_path FROM declaration d JOIN package p ON p.id = d.package_id
		WHERE NOT EXISTS (SELECT 1 FROM symbol s WHERE s.declaration_id = d.id) ORDER BY 1`, func(r *sql.Rows) error {
		var id int64
		var path string
		if err := r.Scan(&id, &path); err != nil {
			return err
		}
		add("declaration %d of %s declares no symbol", id, path)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("declarations: %w", err)
	}

	var want, entries, indexed int
	if err := d.QueryRowContext(ctx, `SELECT (SELECT count(*) FROM package) + (SELECT count(*) FROM symbol) + (SELECT count(*) FROM article_section),
		(SELECT count(*) FROM search_entry), (SELECT count(*) FROM search_docsize)`).Scan(&want, &entries, &indexed); err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	if entries != want {
		add("search_entry holds %d rows, want one for each package, symbol and section: %d", entries, want)
	}
	if indexed != entries {
		add("the search index holds %d rows, not search_entry's %d", indexed, entries)
	}

	if err := checkLinks(ctx, d, add); err != nil {
		return nil, fmt.Errorf("links: %w", err)
	}
	return p, nil
}

// checkLinks reports every link of the reference's Markdown that does not lead to a package,
// symbol, article or section in it.
func checkLinks(ctx context.Context, d Querier, add func(string, ...any)) error {
	packages := map[string]bool{}
	symbols := map[string]string{} // path.Name or path.Type.Name → kind
	anchors := map[string]bool{}   // slug, and slug#anchor
	for _, q := range []struct {
		query string
		fn    func(a, b string)
	}{
		{`SELECT import_path, '' FROM package`, func(a, _ string) { packages[a] = true }},
		{`SELECT p.import_path || '.' || coalesce(s.receiver || '.', '') || s.name, s.kind FROM symbol s JOIN package p ON p.id = s.package_id`,
			func(a, b string) { symbols[a] = b }},
		{`SELECT slug, '' FROM article UNION ALL SELECT article_slug || '#' || anchor, '' FROM article_section`, func(a, _ string) { anchors[a] = true }},
	} {
		if err := eachRow(ctx, d, q.query, func(r *sql.Rows) error {
			var a, b string
			if err := r.Scan(&a, &b); err != nil {
				return err
			}
			q.fn(a, b)
			return nil
		}); err != nil {
			return err
		}
	}
	resolves := func(target string) bool {
		if rest, ok := strings.CutPrefix(target, "godoc-article:"); ok {
			slug, anchor, hasAnchor := strings.Cut(rest, "#")
			return anchors[slug] && (!hasAnchor || anchors[slug+"#"+anchor])
		}
		rest, ok := strings.CutPrefix(target, "godoc:")
		if !ok {
			return false
		}
		path, names := SplitGodoc(rest)
		switch {
		case !packages[path]:
			return false
		case len(names) == 0:
			return true
		case len(names) == 1:
			return symbols[path+"."+names[0]] != ""
		case len(names) == 2: // a method, or a field or interface method of the type
			return symbols[rest] == "method" || symbols[path+"."+names[0]] == "type"
		}
		return false
	}
	return eachRow(ctx, d, `SELECT 'package ' || import_path, doc FROM package WHERE doc IS NOT NULL
		UNION ALL SELECT 'the declaration of ' || (SELECT p.import_path || '.' || coalesce(s.receiver || '.', '') || s.name
			FROM symbol s JOIN package p ON p.id = s.package_id WHERE s.declaration_id = d.id ORDER BY s.position LIMIT 1), d.doc
			FROM declaration d WHERE d.doc IS NOT NULL
		UNION ALL SELECT 'example ' || e.position || ' of ' || p.import_path, e.doc FROM example e JOIN package p ON p.id = e.package_id WHERE e.doc IS NOT NULL
		UNION ALL SELECT 'the lead of article ' || slug, lead FROM article WHERE lead IS NOT NULL
		UNION ALL SELECT 'section ' || anchor || ' of article ' || article_slug, body FROM article_section WHERE body IS NOT NULL`,
		func(r *sql.Rows) error {
			var where sql.NullString
			var text string
			if err := r.Scan(&where, &text); err != nil {
				return err
			}
			for _, m := range reLink.FindAllStringSubmatch(text, -1) {
				if !resolves(m[1]) {
					add("%s links to %s, which is not in the reference", where.String, m[1])
				}
			}
			return nil
		})
}

// SplitGodoc splits what follows godoc: in a link into the package's import path and the
// names after it: "io.Reader.Read" is io, [Reader Read]. Standard-library import paths have no
// dots, so the names start at the first dot after the last slash.
func SplitGodoc(s string) (path string, names []string) {
	slash := strings.LastIndexByte(s, '/')
	dot := strings.IndexByte(s[slash+1:], '.')
	if dot < 0 {
		return s, nil
	}
	dot += slash + 1
	return s[:dot], strings.Split(s[dot+1:], ".")
}

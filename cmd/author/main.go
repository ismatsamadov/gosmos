// Command author writes original content into a Gosmos library: courses, chapters, lessons
// from lesson packages, and a course's story from its story package, the folders writers author
// in. It knows no content source; any library in the format docs/library-format.md describes
// works.
//
// Every write is checked before it is made and runs in one transaction, which commits only if
// the library still passes its checks. Run "author help" for the packages' layouts.
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/jobloom/gosmos/db"
	"github.com/jobloom/gosmos/internal/authoring"
	"github.com/jobloom/gosmos/internal/media"
	"github.com/jobloom/gosmos/internal/store"
)

const usage = `usage:
  author course <file> --slug S --title T --kind 'Course'|'Guided Project' --short S --description D [--position N]
  author chapter <file> --course S --slug S --title T --description D [--position N]
  author lesson <file> --course S --chapter S <package folder>
  author story <file> --course S <story folder>
  author export <file> --course S [--chapter S | --story] <out folder>
  author help

course and chapter add or update one by its slug; a position left out keeps the one it has,
or puts a new one last. lesson imports one lesson package, replacing the lesson of the same
course, chapter and slug in place, with its uuid. story imports a course's story package in
place of the whole story the course had; a folder with none of its files clears it. export
writes a course's lessons, or one chapter's, as packages: <out>/<chapter>/<lesson>/, or
<out>/<lesson>/ for one chapter; with --story, it writes the course's story package into
<out> instead.`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}

// run carries out one command and returns the exit status: 0 for success, 1 for a refusal or
// failure, 2 for a usage mistake.
func run(ctx context.Context, args []string, out, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(errOut, usage)
		return 2
	}
	commands := map[string]func(context.Context, []string, io.Writer) error{
		"course": course, "chapter": chapter, "lesson": lesson, "story": story, "export": export,
	}
	if args[0] == "help" {
		fmt.Fprintf(out, "%s\n\n%s\n\n%s\n", usage, authoring.Layout, authoring.StoryLayout)
		return 0
	}
	cmd, ok := commands[args[0]]
	if !ok {
		fmt.Fprintf(errOut, "unknown command %q\n%s\n", args[0], usage)
		return 2
	}
	if err := cmd(ctx, args[1:], out); err != nil {
		var u usageError
		if errors.As(err, &u) {
			fmt.Fprintf(errOut, "author %s: %v\n%s\n", args[0], err, usage)
			return 2
		}
		fmt.Fprintf(errOut, "author %s: %v\n", args[0], err)
		return 1
	}
	return 0
}

type usageError struct{ error }

// parse reads flags and positional arguments in any order, and wants n positional ones.
func parse(fs *flag.FlagSet, args []string, n int) ([]string, error) {
	fs.SetOutput(io.Discard)
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, usageError{err}
		}
		if fs.NArg() == 0 {
			break
		}
		pos = append(pos, fs.Arg(0))
		args = fs.Args()[1:]
	}
	if len(pos) != n {
		return nil, usageError{fmt.Errorf("want %d arguments besides the flags, got %d", n, len(pos))}
	}
	return pos, nil
}

func now() string { return store.UTC(time.Now()) }

// write opens the library and runs fn in one transaction, which commits only if the library
// still passes its checks.
func write(ctx context.Context, file string, fn func(tx *sql.Tx) error) error {
	d, err := db.Open(file)
	if err != nil {
		return err
	}
	defer d.Close()
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(tx); err != nil {
		return err
	}
	problems, err := db.Validate(ctx, tx)
	if err != nil {
		return err
	}
	if err := problems.Err(); err != nil {
		return fmt.Errorf("nothing was written: %w", err)
	}
	return tx.Commit()
}

// noMedia refuses pictures, videos and widgets in text the app shows.
func noMedia(texts map[string]string) error {
	for _, k := range slices.Sorted(maps.Keys(texts)) {
		if _, st := media.Strip(texts[k]); st.Total() > 0 {
			return fmt.Errorf("%s holds pictures, videos or widgets, which the app never shows; leave them out", k)
		}
	}
	return nil
}

// ---- course ----

func course(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("course", flag.ContinueOnError)
	slug := fs.String("slug", "", "")
	title := fs.String("title", "", "")
	kind := fs.String("kind", "", "")
	short := fs.String("short", "", "")
	description := fs.String("description", "", "")
	position := fs.Int64("position", -1, "")
	pos, err := parse(fs, args, 1)
	if err != nil {
		return err
	}
	switch {
	case !authoring.ValidSlug(*slug):
		return usageError{fmt.Errorf("--slug %q is not lowercase letters, digits and single hyphens", *slug)}
	case *title == "" || *short == "" || *description == "":
		return usageError{errors.New("--title, --short and --description are required")}
	case *kind != "Course" && *kind != "Guided Project":
		return usageError{fmt.Errorf("--kind %q is not Course or Guided Project", *kind)}
	}
	if err := noMedia(map[string]string{"--title": *title, "--short": *short, "--description": *description}); err != nil {
		return err
	}
	return write(ctx, pos[0], func(tx *sql.Tx) error {
		var uuid string
		var at int64
		var cur [4]string
		err := tx.QueryRow(`SELECT uuid, position, title, kind, short_description, description FROM course WHERE slug = ?`, *slug).
			Scan(&uuid, &at, &cur[0], &cur[1], &cur[2], &cur[3])
		exists := err == nil
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if !exists {
			uuid = authoring.CourseUUID(*slug)
			if err := tx.QueryRow(`SELECT coalesce(max(position) + 1, 0) FROM course`).Scan(&at); err != nil {
				return err
			}
		}
		if *position >= 0 {
			at = *position
		}
		var holder string
		if err := tx.QueryRow(`SELECT slug FROM course WHERE position = ? AND uuid <> ?`, at, uuid).Scan(&holder); err == nil {
			return fmt.Errorf("position %d is taken by course %s", at, holder)
		}
		if exists && cur == [4]string{*title, *kind, *short, *description} && *position < 0 {
			fmt.Fprintf(out, "course %s: unchanged (%s)\n", *slug, uuid)
			return nil
		}
		if _, err := tx.Exec(`INSERT INTO course (uuid, slug, title, kind, short_description, description, position) VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(uuid) DO UPDATE SET slug = excluded.slug, title = excluded.title, kind = excluded.kind,
				short_description = excluded.short_description, description = excluded.description, position = excluded.position`,
			uuid, *slug, *title, *kind, *short, *description, at); err != nil {
			return err
		}
		fmt.Fprintf(out, "course %s: %s at position %d (%s)\n", *slug, map[bool]string{true: "updated", false: "added"}[exists], at, uuid)
		return nil
	})
}

// ---- chapter ----

func chapter(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("chapter", flag.ContinueOnError)
	courseSlug := fs.String("course", "", "")
	slug := fs.String("slug", "", "")
	title := fs.String("title", "", "")
	description := fs.String("description", "", "")
	position := fs.Int64("position", -1, "")
	pos, err := parse(fs, args, 1)
	if err != nil {
		return err
	}
	switch {
	case !authoring.ValidSlug(*slug):
		return usageError{fmt.Errorf("--slug %q is not lowercase letters, digits and single hyphens", *slug)}
	case *courseSlug == "" || *title == "" || *description == "":
		return usageError{errors.New("--course, --title and --description are required")}
	}
	if err := noMedia(map[string]string{"--title": *title, "--description": *description}); err != nil {
		return err
	}
	return write(ctx, pos[0], func(tx *sql.Tx) error {
		courseUUID, err := courseBySlug(tx, *courseSlug)
		if err != nil {
			return err
		}
		var uuid, curTitle, curDescription string
		var at int64
		err = tx.QueryRow(`SELECT uuid, position, title, description FROM chapter WHERE course_uuid = ? AND slug = ?`, courseUUID, *slug).
			Scan(&uuid, &at, &curTitle, &curDescription)
		exists := err == nil
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if !exists {
			uuid = authoring.ChapterUUID(*courseSlug, *slug)
			if err := tx.QueryRow(`SELECT coalesce(max(position) + 1, 0) FROM chapter WHERE course_uuid = ?`, courseUUID).Scan(&at); err != nil {
				return err
			}
		}
		if *position >= 0 {
			at = *position
		}
		var holder string
		if err := tx.QueryRow(`SELECT slug FROM chapter WHERE course_uuid = ? AND position = ? AND uuid <> ?`, courseUUID, at, uuid).Scan(&holder); err == nil {
			return fmt.Errorf("position %d of course %s is taken by chapter %s", at, *courseSlug, holder)
		}
		if exists && curTitle == *title && curDescription == *description && *position < 0 {
			fmt.Fprintf(out, "chapter %s/%s: unchanged (%s)\n", *courseSlug, *slug, uuid)
			return nil
		}
		if _, err := tx.Exec(`INSERT INTO chapter (uuid, course_uuid, position, slug, title, description) VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(uuid) DO UPDATE SET course_uuid = excluded.course_uuid, position = excluded.position, slug = excluded.slug,
				title = excluded.title, description = excluded.description`,
			uuid, courseUUID, at, *slug, *title, *description); err != nil {
			return err
		}
		fmt.Fprintf(out, "chapter %s/%s: %s at position %d (%s)\n", *courseSlug, *slug, map[bool]string{true: "updated", false: "added"}[exists], at, uuid)
		return nil
	})
}

func courseBySlug(q store.Querier, slug string) (string, error) {
	var uuid string
	if err := q.QueryRow(`SELECT uuid FROM course WHERE slug = ?`, slug).Scan(&uuid); errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("the library has no course %s; add it with author course", slug)
	} else if err != nil {
		return "", err
	}
	return uuid, nil
}

func chapterBySlug(q store.Querier, courseUUID, course, slug string) (string, error) {
	var uuid string
	if err := q.QueryRow(`SELECT uuid FROM chapter WHERE course_uuid = ? AND slug = ?`, courseUUID, slug).Scan(&uuid); errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("course %s has no chapter %s; add it with author chapter", course, slug)
	} else if err != nil {
		return "", err
	}
	return uuid, nil
}

// chapterSlugs returns the slugs of a course's chapters, in the course's order.
func chapterSlugs(q store.Querier, courseUUID string) ([]string, error) {
	rows, err := q.Query(`SELECT slug FROM chapter WHERE course_uuid = ? ORDER BY position`, courseUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		slugs = append(slugs, slug)
	}
	return slugs, rows.Err()
}

// ---- lesson ----

func lesson(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("lesson", flag.ContinueOnError)
	courseSlug := fs.String("course", "", "")
	chapterSlug := fs.String("chapter", "", "")
	pos, err := parse(fs, args, 2)
	if err != nil {
		return err
	}
	if *courseSlug == "" || *chapterSlug == "" {
		return usageError{errors.New("--course and --chapter are required")}
	}
	pkg, err := authoring.Read(pos[1])
	if err != nil {
		return fmt.Errorf("package %s: %w", pos[1], err)
	}
	name := *courseSlug + "/" + *chapterSlug + "/" + pkg.Slug
	return write(ctx, pos[0], func(tx *sql.Tx) error {
		courseUUID, err := courseBySlug(tx, *courseSlug)
		if err != nil {
			return err
		}
		chapterUUID, err := chapterBySlug(tx, courseUUID, *courseSlug, *chapterSlug)
		if err != nil {
			return err
		}
		var uuid string
		err = tx.QueryRow(`SELECT uuid FROM lesson WHERE chapter_uuid = ? AND slug = ?`, chapterUUID, pkg.Slug).Scan(&uuid)
		exists := err == nil
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if !exists {
			uuid = authoring.LessonUUID(*courseSlug, *chapterSlug, pkg.Slug)
		}
		var holder string
		if err := tx.QueryRow(`SELECT slug FROM lesson WHERE chapter_uuid = ? AND position = ? AND uuid <> ?`, chapterUUID, pkg.Position, uuid).Scan(&holder); err == nil {
			return fmt.Errorf("position %d of chapter %s is taken by lesson %s", pkg.Position, *chapterSlug, holder)
		}
		t := now()
		times := authoring.Times{Created: t, Modified: t}
		var diff []string
		if exists {
			cur, _, err := store.Get(tx, uuid)
			if err != nil {
				return err
			}
			times = authoring.Times{Created: cur.CreatedAt, Modified: cur.LastMod}
			rec := pkg.Record(uuid, chapterUUID, times)
			parts, err := store.ReadParts(tx, uuid)
			if err != nil {
				return err
			}
			if diff = append(store.LessonDiff(cur, rec.Lesson), store.PartsDiff(parts, rec.Rows)...); len(diff) == 0 {
				fmt.Fprintf(out, "lesson %s: unchanged (%s)\n", name, uuid)
				return nil
			}
			times.Modified = t
		}
		if err := store.Put(tx, pkg.Record(uuid, chapterUUID, times)); err != nil {
			return err
		}
		if exists {
			fmt.Fprintf(out, "lesson %s: updated %s (%s)\n", name, strings.Join(diff, ", "), uuid)
		} else {
			fmt.Fprintf(out, "lesson %s: added at position %d (%s)\n", name, pkg.Position, uuid)
		}
		return nil
	})
}

// ---- story ----

func story(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("story", flag.ContinueOnError)
	courseSlug := fs.String("course", "", "")
	pos, err := parse(fs, args, 2)
	if err != nil {
		return err
	}
	if *courseSlug == "" {
		return usageError{errors.New("--course is required")}
	}
	s, err := authoring.ReadStory(pos[1])
	if err != nil {
		return fmt.Errorf("package %s: %w", pos[1], err)
	}
	return write(ctx, pos[0], func(tx *sql.Tx) error {
		courseUUID, err := courseBySlug(tx, *courseSlug)
		if err != nil {
			return err
		}
		chapters, err := chapterSlugs(tx, courseUUID)
		if err != nil {
			return err
		}
		if err := s.CheckChapters(*courseSlug, chapters); err != nil {
			return fmt.Errorf("package %s: %w", pos[1], err)
		}
		cur, err := authoring.LoadStory(tx, courseUUID)
		if err != nil {
			return err
		}
		changed := cur.Changed(s)
		if len(changed) == 0 {
			fmt.Fprintf(out, "story %s: unchanged (%s)\n", *courseSlug, s.Summary())
			return nil
		}
		if err := s.Put(tx, courseUUID); err != nil {
			return err
		}
		switch {
		case cur.Empty():
			fmt.Fprintf(out, "story %s: added (%s)\n", *courseSlug, s.Summary())
		case s.Empty():
			fmt.Fprintf(out, "story %s: cleared (it had %s)\n", *courseSlug, cur.Summary())
		default:
			fmt.Fprintf(out, "story %s: updated %s (%s)\n", *courseSlug, strings.Join(changed, ", "), s.Summary())
		}
		return nil
	})
}

// ---- export ----

func export(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	courseSlug := fs.String("course", "", "")
	chapterSlug := fs.String("chapter", "", "")
	wantStory := fs.Bool("story", false, "")
	pos, err := parse(fs, args, 2)
	if err != nil {
		return err
	}
	switch {
	case *courseSlug == "":
		return usageError{errors.New("--course is required")}
	case *wantStory && *chapterSlug != "":
		return usageError{errors.New("--chapter and --story do not go together: the story is the whole course's")}
	}
	file, outDir := pos[0], pos[1]
	if entries, err := os.ReadDir(outDir); err == nil && len(entries) > 0 {
		return fmt.Errorf("%s is not empty; export into a new folder", outDir)
	}
	d, recovered, err := db.Inspect(file)
	if err != nil {
		return err
	}
	defer d.Close()
	if recovered {
		fmt.Fprintln(out, db.RecoveredMessage(file))
	}
	var version int
	if err := d.QueryRowContext(ctx, `PRAGMA user_version`).Scan(&version); err != nil {
		return err
	}
	if version != db.Version {
		return fmt.Errorf("%s is of format %d, not %d; migrate it first (library migrate)", file, version, db.Version)
	}
	courseUUID, err := courseBySlug(d, *courseSlug)
	if err != nil {
		return err
	}
	if *wantStory {
		s, err := authoring.LoadStory(d, courseUUID)
		if err != nil {
			return err
		}
		if err := s.Write(outDir); err != nil {
			return err
		}
		fmt.Fprintf(out, "exported the story of %s to %s (%s)\n", *courseSlug, outDir, s.Summary())
		return nil
	}
	query := `SELECT l.uuid, ch.slug FROM lesson l JOIN chapter ch ON ch.uuid = l.chapter_uuid WHERE ch.course_uuid = ?`
	qargs := []any{courseUUID}
	if *chapterSlug != "" {
		if _, err := chapterBySlug(d, courseUUID, *courseSlug, *chapterSlug); err != nil {
			return err
		}
		query += ` AND ch.slug = ?`
		qargs = append(qargs, *chapterSlug)
	}
	rows, err := d.QueryContext(ctx, query+` ORDER BY ch.position, l.position`, qargs...)
	if err != nil {
		return err
	}
	type item struct{ uuid, chapter string }
	var items []item
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.uuid, &it.chapter); err != nil {
			rows.Close()
			return err
		}
		items = append(items, it)
	}
	rows.Close()
	written := 0
	for _, it := range items {
		pkg, err := authoring.Load(d, it.uuid)
		if err != nil {
			return err
		}
		dir := filepath.Join(outDir, it.chapter, pkg.Slug)
		if *chapterSlug != "" {
			dir = filepath.Join(outDir, pkg.Slug)
		}
		if err := pkg.Write(dir); err != nil {
			return err
		}
		written++
	}
	fmt.Fprintf(out, "exported %d lessons of %s to %s\n", written, *courseSlug, outDir)
	return nil
}

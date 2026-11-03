// Command library makes, checks, migrates and verifies Gosmos course libraries: SQLite files in
// the format docs/library-format.md describes, whatever they hold and wherever it came from.
//
//	library create <file> [--title T] [--version V]
//	library check <file>
//	library migrate <file> [--learner <file>] [--title T]
//	library verify <file> [--course <slug>] [--teaching]
//
// check prints one line for each problem and exits 1 if there is any. It writes nothing, with
// one exception: a write that was interrupted (a crash, or Ctrl-C during a write) leaves a
// journal that has to be rolled back before the file can be read, and check does that first and
// says so, as SQLite itself would on any open that may write.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jobloom/gosmos/db"
)

const usage = `usage:
  library create <file> [--title T] [--version V]
      make a new, empty library, titled with the file's name and of version 1 unless given
  library check <file>
      run the checks the app makes before it opens a library
  library migrate <file> [--learner <file>] [--title T]
      bring a library of an older format up to date. Going to format 7, the library gets its
      title (else its first course's, else the file's name), and the learner's data in it
      moves to the learner database, made from learner.sql if it does not exist; a library
      that holds such data needs --learner. Format 8 has code, tests and quiz lessons only: a
      library with command-line, GitHub or interview lessons is refused. Format 9 adds the
      tables of a course's story, empty
  library verify <file> [--course <slug>] [--teaching]
      run every code lesson with the Go toolchain: it must pass with its solution and fail
      with its starter. With --teaching, every code lesson must also have 2 or 3 hints and a
      solution note, and every quiz answer an explanation`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}

// run carries out one command and returns the exit status: 0 for success, 1 for a library
// that fails or a command that could not be done, 2 for a usage mistake.
func run(ctx context.Context, args []string, out, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(errOut, usage)
		return 2
	}
	cmd := args[0]
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var title, version, learner, course string
	var teaching bool
	switch cmd {
	case "create":
		fs.StringVar(&title, "title", "", "")
		fs.StringVar(&version, "version", "1", "")
	case "check":
	case "migrate":
		fs.StringVar(&learner, "learner", "", "")
		fs.StringVar(&title, "title", "", "")
	case "verify":
		fs.StringVar(&course, "course", "", "")
		fs.BoolVar(&teaching, "teaching", false, "")
	default:
		fmt.Fprintf(errOut, "unknown command %q\n%s\n", cmd, usage)
		return 2
	}
	pos, err := parse(fs, args[1:])
	if err == nil && len(pos) != 1 {
		err = fmt.Errorf("want one file, got %d arguments", len(pos))
	}
	if err != nil {
		fmt.Fprintf(errOut, "library %s: %v\n%s\n", cmd, err, usage)
		return 2
	}
	path := pos[0]
	var ok bool
	switch cmd {
	case "create":
		err = create(path, title, version, out)
		ok = err == nil
	case "check":
		ok, err = check(ctx, path, out)
	case "migrate":
		err = migrate(ctx, path, db.MigrateOptions{Title: title, Learner: learner}, out)
		ok = err == nil
	case "verify":
		ok, err = verify(ctx, path, course, teaching, out)
	}
	if err != nil {
		fmt.Fprintln(errOut, "library:", err)
		return 1
	}
	if !ok {
		return 1
	}
	return 0
}

// parse reads flags and positional arguments in any order.
func parse(fs *flag.FlagSet, args []string) ([]string, error) {
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return pos, nil
		}
		pos, args = append(pos, fs.Arg(0)), fs.Args()[1:]
	}
}

// fileTitle is the title a library gets from its file's name: the name without its extension.
func fileTitle(path string) string {
	base := filepath.Base(path)
	if name := strings.TrimSuffix(base, filepath.Ext(base)); name != "" {
		return name
	}
	return base
}

// create makes a new, empty library at path. Nothing may be there yet.
func create(path, title, version string, out io.Writer) error {
	if title = strings.TrimSpace(title); title == "" {
		title = fileTitle(path)
	}
	if version = strings.TrimSpace(version); version == "" {
		return errors.New("--version must not be empty")
	}
	d, err := db.Create(path, db.Info{Title: title, Version: version})
	if err != nil {
		return err
	}
	info, err := db.ReadInfo(context.Background(), d)
	if err := errors.Join(err, d.Close()); err != nil {
		return err
	}
	fmt.Fprintf(out, "created %s: an empty library of format %d, %s\n", path, db.Version, describeInfo(info))
	return nil
}

// describeInfo names a library by its library_info row.
func describeInfo(i db.Info) string {
	return fmt.Sprintf("%q version %s (uuid %s)", i.Title, i.Version, i.UUID)
}

// check runs the library checks on the file as it is, printing one line per problem, and
// reports whether it passed. It writes nothing.
func check(ctx context.Context, path string, out io.Writer) (bool, error) {
	d, recovered, err := db.Inspect(path)
	if err != nil {
		return false, err
	}
	defer d.Close()
	if recovered {
		fmt.Fprintln(out, db.RecoveredMessage(path))
	}
	problems, err := db.Validate(ctx, d)
	if err != nil {
		return false, err
	}
	for _, p := range problems {
		fmt.Fprintln(out, p)
	}
	if len(problems) > 0 {
		return false, nil
	}
	info, err := db.ReadInfo(ctx, d)
	if err != nil {
		return false, err
	}
	var courses, chapters, lessons int
	if err := d.QueryRowContext(ctx, `SELECT (SELECT count(*) FROM course), (SELECT count(*) FROM chapter), (SELECT count(*) FROM lesson)`).
		Scan(&courses, &chapters, &lessons); err != nil {
		return false, err
	}
	fmt.Fprintf(out, "ok: %s is %s, a library of format %d with %s, %s and %s\n",
		path, describeInfo(info), db.Version, plural(courses, "course"), plural(chapters, "chapter"), plural(lessons, "lesson"))
	return true, nil
}

// migrate brings a library of an older format up to date (db.Migrate), says what it did, and
// checks the result.
func migrate(ctx context.Context, path string, opt db.MigrateOptions, out io.Writer) error {
	m, err := db.Migrate(ctx, path, opt)
	if m.Recovered {
		fmt.Fprintln(out, db.RecoveredMessage(path))
	}
	if err != nil {
		return err
	}
	if m.From == m.To {
		fmt.Fprintf(out, "%s is already of format %d: %s\n", path, m.To, describeInfo(m.Info))
		return nil
	}
	fmt.Fprintf(out, "migrated %s from format %d to %d: %s\n", path, m.From, m.To, describeInfo(m.Info))
	if m.Learner != "" {
		var counts []string
		for _, r := range m.Moved {
			counts = append(counts, fmt.Sprintf("%s %d", r.Table, r.Rows))
		}
		switch {
		case m.Already:
			fmt.Fprintf(out, "the learner database %s already held its learner data, from an earlier run, so none was copied again: %s\n", m.Learner, strings.Join(counts, ", "))
		case m.Created:
			fmt.Fprintf(out, "moved its learner data into %s, a new learner database: %s\n", m.Learner, strings.Join(counts, ", "))
		default:
			fmt.Fprintf(out, "moved its learner data into the learner database %s: %s\n", m.Learner, strings.Join(counts, ", "))
		}
	}
	var dropped []string
	for _, r := range m.Dropped {
		if r.Rows > 0 {
			dropped = append(dropped, fmt.Sprintf("%s %d", r.Table, r.Rows))
		}
	}
	if len(dropped) > 0 {
		fmt.Fprintf(out, "dropped the learner data only command-line lessons used, which Gosmos no longer has: %s\n", strings.Join(dropped, ", "))
	}
	ok, err := check(ctx, path, out)
	if err == nil && !ok {
		err = errors.New("the migrated library fails its checks")
	}
	return err
}

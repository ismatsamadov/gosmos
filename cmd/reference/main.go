// Command reference makes and checks the Go reference built into Gosmos, reference/go.sqlite:
// the standard library's documentation and five articles from go.dev, in the format
// db/reference.sql describes, so that the app needs no network to look anything up.
//
//	reference build --goroot <dir> --website <dir> --website-commit <sha> --out <file>
//	reference check <file>
//
// build reads the standard library of the Go at --goroot (go list std without internal and
// vendored packages, for darwin/arm64) with go/parser and go/doc, its doc/go_spec.html and
// doc/go_mem.html, and Effective Go, the FAQ and How to Write Go Code from --website, a checkout
// of golang/website at --website-commit, or a folder holding its files under _content/ as the
// checkout does. It uses no network. check prints one line for each problem and exits 1 if
// there is any; it writes nothing.
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jobloom/gosmos/db"
)

const usage = `usage:
  reference build --goroot <dir> --website <dir> --website-commit <sha> --out <file>
      make the Go reference anew from the Go at goroot and golang/website's files at that
      commit; the file at --out is replaced only once the new one passes its checks
  reference check <file>
      run the reference's checks, and say what it holds`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}

// run carries out one command and returns the exit status: 0 for success, 1 for a reference
// that fails or a command that could not be done, 2 for a usage mistake.
func run(ctx context.Context, args []string, out, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(errOut, usage)
		return 2
	}
	cmd := args[0]
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var goroot, website, commit, file string
	wantArgs := 0
	switch cmd {
	case "build":
		fs.StringVar(&goroot, "goroot", "", "")
		fs.StringVar(&website, "website", "", "")
		fs.StringVar(&commit, "website-commit", "", "")
		fs.StringVar(&file, "out", "", "")
	case "check":
		wantArgs = 1
	default:
		fmt.Fprintf(errOut, "unknown command %q\n%s\n", cmd, usage)
		return 2
	}
	pos, err := parse(fs, args[1:])
	switch {
	case err != nil:
	case len(pos) != wantArgs:
		err = fmt.Errorf("want %d file arguments, got %d", wantArgs, len(pos))
	case cmd == "build" && (goroot == "" || website == "" || commit == "" || file == ""):
		err = errors.New("--goroot, --website, --website-commit and --out are all needed")
	}
	if err != nil {
		fmt.Fprintf(errOut, "reference %s: %v\n%s\n", cmd, err, usage)
		return 2
	}
	var ok bool
	switch cmd {
	case "build":
		var c counts
		if c, err = buildReference(ctx, goroot, website, commit, file, time.Now()); err == nil {
			ok = true
			size := ""
			if info, err := os.Stat(file); err == nil {
				size = fmt.Sprintf(" (%.1f MB)", float64(info.Size())/1e6)
			}
			fmt.Fprintf(out, "built %s%s: %s\n", file, size, c)
		}
	case "check":
		ok, err = check(ctx, pos[0], out)
	}
	if err != nil {
		fmt.Fprintln(errOut, "reference:", err)
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

// check runs the reference's checks on the file, printing one line per problem, and reports
// whether it passed.
func check(ctx context.Context, path string, out io.Writer) (bool, error) {
	d, err := db.OpenReference(path)
	if err != nil {
		return false, err
	}
	defer d.Close()
	problems, err := db.ValidateReference(ctx, d)
	if err != nil {
		return false, err
	}
	for _, p := range problems {
		fmt.Fprintln(out, p)
	}
	if len(problems) > 0 {
		return false, nil
	}
	var version, platform, commit string
	var c counts
	if err := d.QueryRowContext(ctx, `SELECT go_version, platform, website_commit,
		(SELECT count(*) FROM package), (SELECT count(*) FROM symbol), (SELECT count(*) FROM example),
		(SELECT count(*) FROM article), (SELECT count(*) FROM article_section) FROM source`).
		Scan(&version, &platform, &commit, &c.packages, &c.symbols, &c.examples, &c.articles, &c.sections); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("the reference has no source row")
		}
		return false, err
	}
	fmt.Fprintf(out, "ok: %s is the Go reference for %s on %s (golang/website %.12s), with %s\n", path, version, platform, commit, c)
	return true, nil
}

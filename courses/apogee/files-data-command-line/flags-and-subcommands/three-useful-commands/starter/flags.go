package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
)

// newFlagSet makes a command's flag set, whose usage starts with the command's synopsis.
func newFlagSet(name, synopsis string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "usage: gctl %s\n", synopsis)
		fs.PrintDefaults()
	}
	return fs
}

// parseFlags parses args with fs and reports whether the command goes on; if not, status is its exit status.
func parseFlags(fs *flag.FlagSet, args []string, stdout, stderr io.Writer) (status int, ok bool) {
	var messages bytes.Buffer
	fs.SetOutput(&messages)
	err := fs.Parse(args)
	if errors.Is(err, flag.ErrHelp) {
		messages.WriteTo(stdout)
		return 0, false
	}
	if err != nil {
		messages.WriteTo(stderr)
		return 2, false
	}
	return 0, true
}

// usageError writes a usage mistake and the command's usage to stderr and returns status 2.
func usageError(fs *flag.FlagSet, stderr io.Writer, problem string) int {
	fmt.Fprintln(stderr, problem)
	fs.SetOutput(stderr)
	fs.Usage()
	return 2
}

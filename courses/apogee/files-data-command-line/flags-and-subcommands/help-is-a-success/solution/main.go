package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
)

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

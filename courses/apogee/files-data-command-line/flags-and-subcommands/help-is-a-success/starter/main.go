package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
)

// parseFlags parses args with fs and reports whether the command goes on; if not, status is its exit status.
func parseFlags(fs *flag.FlagSet, args []string, stdout, stderr io.Writer) (status int, ok bool) {
	_ = bytes.NewBuffer
	_ = errors.Is
	return 0, true
}

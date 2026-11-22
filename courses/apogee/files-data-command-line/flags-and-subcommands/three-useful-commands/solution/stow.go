package main

import (
	"fmt"
	"io"
)

// runStow records that a station's dishes are stowed, and why.
func runStow(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("stow", "stow [-reason TEXT] CODE")
	reason := fs.String("reason", "no reason given", "record `TEXT` as the reason for the stow")
	if status, ok := parseFlags(fs, args, stdout, stderr); !ok {
		return status
	}
	if fs.NArg() != 1 {
		return usageError(fs, stderr, "stow needs one station code, after its flags")
	}
	code := fs.Arg(0)
	if _, ok := findStation(code); !ok {
		return report(nil, fmt.Errorf("unknown station %q", code), stdout, stderr)
	}
	return report([]string{code + " stowed: " + *reason}, nil, stdout, stderr)
}

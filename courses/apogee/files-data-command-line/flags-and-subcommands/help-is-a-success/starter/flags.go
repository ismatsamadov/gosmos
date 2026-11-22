package main

import (
	"flag"
	"fmt"
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

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// run is gctl: it runs the command in args and returns its exit status.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, usageText, err := parseArgs(args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(stdout, usageText)
		return 0
	}
	if err != nil {
		fmt.Fprintf(stderr, "%v\n%s", err, usageText)
		return 2
	}
	plan, err := planDay(opts.Dir, opts.Format)
	if err == nil {
		err = writeReport(stdout, plan)
	}
	if err != nil {
		fmt.Fprintf(stderr, "gctl: %v\n", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

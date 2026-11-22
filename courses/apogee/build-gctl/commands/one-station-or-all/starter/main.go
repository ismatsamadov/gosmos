package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run starts gctl's shell on stdin and stdout and returns the exit status: 0 when the
// session ends, 1 when reading stdin fails, 2 when the shell is given arguments.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintln(stderr, "usage: gctl")
		return 2
	}
	if err := newShell().Run(stdin, stdout); err != nil {
		fmt.Fprintf(stderr, "gctl: %v\n", err)
		return 1
	}
	return 0
}

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
	_ = errors.Is
	_ = flag.ErrHelp
	_ = fmt.Fprint
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

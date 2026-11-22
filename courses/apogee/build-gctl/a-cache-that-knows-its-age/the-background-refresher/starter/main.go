package main

import (
	"fmt"
	"io"
	"os"
	"time"
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
	s := newShell()
	go refreshEvery(s.cache, time.Minute)
	if err := s.Run(stdin, stdout); err != nil {
		fmt.Fprintf(stderr, "gctl: %v\n", err)
		return 1
	}
	return 0
}

// refreshEvery loads the cache's passes again every interval, for as long as gctl runs. A
// failed load keeps the old passes, and passes tries again once they are stale.
func refreshEvery(c *passCache, interval time.Duration) {
	for range time.Tick(interval) {
		c.refresh()
	}
}

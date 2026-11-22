package main

import (
	"fmt"
	"io"
	"os"
)

// usage is gctl's help: its commands and their arguments.
const usage = `usage: gctl <command> [flags] [arguments]

commands:
  passes [-station CODE] [file]  list the bookings by AOS
  stations                       print the station table
  stow [-reason TEXT] CODE       record a stow
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run hands the arguments after a command's name to that command and returns its exit status.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, "gctl needs a command\n"+usage)
		return 2
	}
	switch args[0] {
	case "-h", "-help", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	case "passes":
		return runPasses(args[1:], stdin, stdout, stderr)
	case "stations":
		return runStations(args[1:], stdout, stderr)
	case "stow":
		return runStow(args[1:], stdout, stderr)
	}
	fmt.Fprintf(stderr, "unknown command %q\n%s", args[0], usage)
	return 2
}

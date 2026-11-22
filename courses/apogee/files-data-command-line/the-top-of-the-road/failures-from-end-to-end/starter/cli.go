package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// usage is gctl's help: its commands and their arguments. This copy of gctl carries only its day command.
const usage = `usage: gctl <command> [flags] [arguments]

commands:
  day [-format json|csv] dir  run a calibration day from dir
`

// dayUsage is the day command's own usage.
const dayUsage = "usage: gctl day [-format json|csv] dir\n"

// options are what one run of gctl day was asked for.
type options struct {
	Dir    string
	Format string
}

// parseArgs reads gctl's arguments and returns the usage that fits them: gctl's own, or its day command's.
// Its error is flag.ErrHelp when help was asked for, and says what is wrong for any other misuse.
func parseArgs(args []string) (options, string, error) {
	if len(args) == 0 {
		return options{}, usage, errors.New("gctl needs a command")
	}
	if args[0] == "-h" || args[0] == "-help" || args[0] == "--help" {
		return options{}, usage, flag.ErrHelp
	}
	if args[0] != "day" {
		return options{}, usage, fmt.Errorf("unknown command %q", args[0])
	}
	var opts options
	flags := flag.NewFlagSet("day", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.Format, "format", "json", "the bookings file's format: json or csv")
	if err := flags.Parse(args[1:]); err != nil {
		return options{}, dayUsage, err
	}
	if opts.Format != "json" && opts.Format != "csv" {
		return options{}, dayUsage, fmt.Errorf("-format %q: want json or csv", opts.Format)
	}
	if flags.NArg() != 1 {
		return options{}, dayUsage, errors.New("day needs one folder, after its flags")
	}
	opts.Dir = flags.Arg(0)
	return opts, dayUsage, nil
}

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

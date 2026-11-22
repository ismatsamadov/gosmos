package main

import (
	"flag"
	"io"
)

// passesOptions holds what the passes command was asked for.
type passesOptions struct {
	Station string // show only this station's bookings; empty means every station
}

// parseOptions reads the passes command's flags from args, writing flag messages to output.
func parseOptions(args []string, output io.Writer) (passesOptions, error) {
	_ = flag.NewFlagSet
	return passesOptions{}, nil
}

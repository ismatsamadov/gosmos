package main

import (
	"fmt"
	"io"
	"maps"
	"slices"
)

// command is one entry of gctl's command table: its help text and the function that runs it.
type command struct {
	help string
	run  func(s *Shell, args []string, out io.Writer) error
}

// commandTable returns every command gctl knows, by name.
func commandTable() map[string]command {
	return map[string]command{
		"help": {help: "list every command", run: cmdHelp},
	}
}

// cmdHelp lists every command and its help text, sorted by name.
func cmdHelp(s *Shell, args []string, out io.Writer) error {
	for _, name := range slices.Sorted(maps.Keys(s.commands)) {
		fmt.Fprintf(out, "%-8s  %s\n", name, s.commands[name].help)
	}
	return nil
}

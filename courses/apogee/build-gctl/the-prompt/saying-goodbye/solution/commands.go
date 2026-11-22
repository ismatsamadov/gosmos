package main

import (
	"errors"
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

// errExit is what the exit command returns to end the session.
var errExit = errors.New("exit")

// commandTable returns every command gctl knows, by name.
func commandTable() map[string]command {
	return map[string]command{
		"exit": {help: "leave the shell", run: cmdExit},
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

// cmdExit ends the session.
func cmdExit(s *Shell, args []string, out io.Writer) error {
	return errExit
}

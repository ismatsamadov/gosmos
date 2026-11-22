package main

import "io"

// command is one entry of gctl's command table: its help text and the function that runs it.
type command struct {
	help string
	run  func(s *Shell, args []string, out io.Writer) error
}

// commandTable returns every command gctl knows, by name.
func commandTable() map[string]command {
	return map[string]command{}
}

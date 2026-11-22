package main

import (
	"bufio"
	"fmt"
	"io"
)

// Shell is gctl's interactive prompt: it reads commands and answers them.
type Shell struct {
	prompt   string
	commands map[string]command
}

// newShell returns a shell with gctl's prompt and command table.
func newShell() *Shell {
	return &Shell{prompt: "gctl> ", commands: commandTable()}
}

// Run reads commands from in, one per line, and writes the prompt and every answer to out.
func (s *Shell) Run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	fmt.Fprint(out, s.prompt)
	for sc.Scan() {
		s.dispatch(sc.Text(), out)
		fmt.Fprint(out, s.prompt)
	}
	return nil
}

// dispatch answers one line typed at the prompt, and returns the error of any command it runs.
func (s *Shell) dispatch(line string, out io.Writer) error {
	fmt.Fprintf(out, "unknown command %q; type help\n", line)
	return nil
}

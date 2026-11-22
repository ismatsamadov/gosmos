package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Shell is gctl's interactive prompt: it reads commands and answers them.
type Shell struct {
	prompt   string
	commands map[string]command
	passes   []Pass
}

// newShell returns a shell with gctl's prompt, its command table and tonight's passes.
func newShell() *Shell {
	return &Shell{prompt: "gctl> ", commands: commandTable(), passes: tonight()}
}

// Run reads commands from in, one per line, and writes the prompt and every answer to out.
// It returns only a failure to read in.
func (s *Shell) Run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	fmt.Fprint(out, s.prompt)
	for sc.Scan() {
		err := s.dispatch(sc.Text(), out)
		if errors.Is(err, errExit) {
			break
		}
		if err != nil {
			fmt.Fprintf(out, "error: %v\n", err)
		}
		fmt.Fprint(out, s.prompt)
	}
	if err := sc.Err(); err != nil {
		return err
	}
	fmt.Fprintln(out, "bye")
	return nil
}

// dispatch answers one line typed at the prompt, and returns the error of any command it runs.
func (s *Shell) dispatch(line string, out io.Writer) error {
	words := strings.Fields(line)
	if len(words) == 0 {
		return nil
	}
	c, ok := s.commands[words[0]]
	if !ok {
		fmt.Fprintf(out, "unknown command %q; type help\n", words[0])
		return nil
	}
	return c.run(s, words[1:], out)
}

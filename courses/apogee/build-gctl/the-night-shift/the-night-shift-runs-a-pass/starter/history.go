package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// history keeps the commands typed at the prompt, and the file they are saved in.
type history struct {
	lines  []string
	path   string // the history file; "" keeps the history in memory only
	failed bool   // a file error was reported this session
}

// add records one typed line before it runs, and appends it to the history file if there is
// one; blank lines are not recorded. It returns a file error the first time one happens.
func (h *history) add(line string) error {
	if strings.TrimSpace(line) == "" {
		return nil
	}
	h.lines = append(h.lines, line)
	if h.path == "" || h.failed {
		return nil
	}
	if err := appendLine(h.path, line); err != nil {
		h.failed = true
		return fmt.Errorf("history: %w", err)
	}
	return nil
}

// load reads the history file into h, if there is one; a missing file is an empty history.
func (h *history) load() error {
	if h.path == "" {
		return nil
	}
	data, err := os.ReadFile(h.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		h.failed = true
		return fmt.Errorf("history: %w", err)
	}
	for line := range strings.Lines(string(data)) {
		h.lines = append(h.lines, strings.TrimSuffix(line, "\n"))
	}
	return nil
}

// appendLine adds line to the end of the file at path, making the file if it is missing.
func appendLine(path, line string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, line)
	return errors.Join(err, f.Close())
}

// cmdHistory lists every command typed so far, numbered from 1.
func cmdHistory(s *Shell, args []string, out io.Writer) error {
	if len(args) > 0 {
		return errors.New("history takes no arguments")
	}
	for i, line := range s.history.lines {
		fmt.Fprintf(out, "%3d  %s\n", i+1, line)
	}
	return nil
}

package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// history keeps the commands typed at the prompt, and the file they are saved in.
type history struct {
	lines  []string
	path   string // the history file; "" keeps the history in memory only
	failed bool   // a file error was reported this session
}

// add records one typed line before it runs; blank lines are not recorded.
func (h *history) add(line string) {
	if strings.TrimSpace(line) != "" {
		h.lines = append(h.lines, line)
	}
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

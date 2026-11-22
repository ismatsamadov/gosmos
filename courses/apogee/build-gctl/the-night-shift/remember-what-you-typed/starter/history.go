package main

import "io"

// history keeps the commands typed at the prompt, and the file they are saved in.
type history struct {
	lines  []string
	path   string // the history file; "" keeps the history in memory only
	failed bool   // a file error was reported this session
}

// add records one typed line before it runs; blank lines are not recorded.
func (h *history) add(line string) {
}

// cmdHistory lists every command typed so far, numbered from 1.
func cmdHistory(s *Shell, args []string, out io.Writer) error {
	return nil
}

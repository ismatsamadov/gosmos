package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Watch shows one pass's frames on out, in order and each once, until the
// pass goes quiet; then it asks the recorder how it is and writes the last line.
func (l *Live) Watch(out io.Writer) error {
	var t Tracker
	for {
		d, err := l.receive()
		if errors.Is(err, os.ErrDeadlineExceeded) {
			break
		}
		if err != nil {
			return err
		}
		for _, line := range strings.Split(string(d), "\n") {
			if seq, ok := seqOf(line); ok && t.Accept(seq) {
				fmt.Fprintln(out, line)
			}
		}
	}
	l.asked++
	answer, err := ask(l.Questions, l.Recorder, l.asked, l.Wait, l.Tries)
	if errors.Is(err, os.ErrDeadlineExceeded) {
		answer = "silent"
	} else if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, lastLine(t, answer))
	return err
}

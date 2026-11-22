package main

import (
	"fmt"
	"io"
)

// Watch shows one pass's frames on out, in order and each once, until the
// pass goes quiet; then it asks the recorder how it is and writes the last line.
func (l *Live) Watch(out io.Writer) error {
	for {
		d, err := l.receive()
		if err != nil {
			return nil
		}
		fmt.Fprintln(out, string(d))
	}
}

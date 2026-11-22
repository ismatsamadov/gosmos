package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
)

// ReadPasses yields the passes of a pass archive, one a line. At a
// line it can't read, or when reading fails, it yields the error
// and stops.
func ReadPasses(r io.Reader) iter.Seq2[PassView, error] {
	return func(yield func(PassView, error) bool) {
		sc := bufio.NewScanner(r)
		for n := 1; sc.Scan(); n++ {
			p, err := parsePass(sc.Text())
			if err != nil {
				err = fmt.Errorf("archive line %d: %w", n, err)
				yield(PassView{}, err)
				return
			}
			if !yield(p, nil) {
				return
			}
		}
		if err := sc.Err(); err != nil {
			yield(PassView{}, err)
		}
	}
}

package main

import (
	"errors"
	"io"
	"strings"
	"testing"
)

const (
	note    = "# night shift 2027-07-06, Paz\n"
	reading = "2027-07-06T03:14:07Z\tLG\t183.25\t41.70\n"
)

func TestCountLines(t *testing.T) {
	for _, tc := range []struct {
		name, log string
		want      int
	}{
		{"two readings", reading + reading, 2},
		{"a note and a blank line", note + reading + "\n" + reading, 4},
		{"no final newline", reading + strings.TrimSuffix(reading, "\n"), 2},
		{"one blank line", "\n", 1},
		{"empty log", "", 0},
		{"the whole log", strings.Repeat(reading, 200000), 200000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := countLines(strings.NewReader(tc.log))
			if got != tc.want || err != nil {
				t.Errorf("countLines(%s) = %d, %v, want %d, nil", tc.name, got, err, tc.want)
			}
		})
	}
	problem := errors.New("stick unplugged")
	damaged := io.MultiReader(strings.NewReader(reading), readerFunc(func([]byte) (int, error) { return 0, problem }))
	got, err := countLines(damaged)
	if got != 1 || !errors.Is(err, problem) {
		t.Errorf("countLines(damaged copy) = %d, %v, want 1, %v", got, err, problem)
	}
}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the supplied test function.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

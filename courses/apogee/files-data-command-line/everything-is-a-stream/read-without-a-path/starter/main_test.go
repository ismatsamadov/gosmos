package main

import (
	"errors"
	"strings"
	"testing"
)

func TestReadText(t *testing.T) {
	for _, want := range []string{"PEL ready\n", "", "  exact  "} {
		got, err := readText(strings.NewReader(want))
		if got != want || err != nil {
			t.Errorf("readText(reader) = %q, %v, want %q, nil", got, err, want)
		}
	}
	problem := errors.New("recording stopped")
	got, err := readText(readerFunc(func(p []byte) (int, error) { return copy(p, "partial"), problem }))
	if got != "partial" || !errors.Is(err, problem) {
		t.Errorf("readText(failing) = %q, %v, want partial and original error", got, err)
	}
}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the supplied test reader.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

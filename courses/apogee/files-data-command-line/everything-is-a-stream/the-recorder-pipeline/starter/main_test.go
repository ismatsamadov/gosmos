package main

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRecordPreview(t *testing.T) {
	var a, b bytes.Buffer
	src := strings.NewReader("abcdef")
	n, err := recordPreview(src, 3, []io.Writer{&a, &b})
	rest, _ := io.ReadAll(src)
	if n != 3 || a.String() != "abc" || b.String() != "abc" || string(rest) != "def" || err != nil {
		t.Errorf("recordPreview() = %d, %q/%q, rest %q, %v, want bounded identical copies and nil", n, a.String(), b.String(), rest, err)
	}
	for _, tc := range []struct {
		limit   int64
		targets []io.Writer
	}{{-1, []io.Writer{&a}}, {3, nil}} {
		src := strings.NewReader("untouched")
		_, err := recordPreview(src, tc.limit, tc.targets)
		if err == nil || src.Len() != 9 {
			t.Errorf("recordPreview(invalid) error = %v, source left %d, want error and 9", err, src.Len())
		}
	}
	a.Reset()
	n, err = recordPreview(strings.NewReader("abc"), 0, []io.Writer{&a})
	if n != 0 || a.Len() != 0 || err != nil {
		t.Errorf("recordPreview(zero) = %d, %v, want 0, nil", n, err)
	}
	problem := errors.New("archive full")
	_, err = recordPreview(strings.NewReader("abc"), 3, []io.Writer{writerFunc(func(p []byte) (int, error) { return 0, problem })})
	if !errors.Is(err, problem) {
		t.Errorf("recordPreview(failed target) error = %v, want original error", err)
	}
}

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the supplied test writer.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

package main

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestPreview(t *testing.T) {
	for _, tc := range []struct {
		name       string
		n          int64
		want, rest string
	}{{"prefix", 3, "abc", "def"}, {"zero", 0, "", "abcdef"}, {"long", 99, "abcdef", ""}, {"negative", -1, "", "abcdef"}} {
		t.Run(tc.name, func(t *testing.T) {
			src := strings.NewReader("abcdef")
			got, err := preview(src, tc.n)
			rest, readErr := io.ReadAll(src)
			if got != tc.want || string(rest) != tc.rest || err != nil || readErr != nil {
				t.Errorf("preview(%d) = %q with rest %q, errors %v/%v, want %q with rest %q", tc.n, got, rest, err, readErr, tc.want, tc.rest)
			}
		})
	}
	problem := errors.New("preview source stopped")
	_, err := preview(readerFunc(func(p []byte) (int, error) { return 0, problem }), 3)
	if !errors.Is(err, problem) {
		t.Errorf("preview(failing source) error = %v, want original error", err)
	}

}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the supplied test function.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

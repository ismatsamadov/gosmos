package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestCopyRecording(t *testing.T) {
	for _, text := range []string{"recording\n", ""} {
		var b bytes.Buffer
		n, err := copyRecording(&b, strings.NewReader(text))
		if b.String() != text || n != int64(len(text)) || err != nil {
			t.Errorf("copyRecording() = %q, %d, %v, want %q, %d, nil", b.String(), n, err, text, len(text))
		}
	}
	problem := errors.New("input lost")
	var b bytes.Buffer
	n, err := copyRecording(&b, readerFunc(func(p []byte) (int, error) { return copy(p, "part"), problem }))
	if n != 4 || b.String() != "part" || !errors.Is(err, problem) {
		t.Errorf("copyRecording(read error) = %d, %q, %v, want 4, part, original error", n, b.String(), err)
	}
	n, err = copyRecording(writerFunc(func(p []byte) (int, error) { return 0, problem }), strings.NewReader("x"))
	if n != 0 || !errors.Is(err, problem) {
		t.Errorf("copyRecording(write error) = %d, %v, want 0, original error", n, err)
	}
}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the supplied test reader.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the supplied test writer.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

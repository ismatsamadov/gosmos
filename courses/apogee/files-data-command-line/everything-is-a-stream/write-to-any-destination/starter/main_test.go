package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriteStatus(t *testing.T) {
	var b bytes.Buffer
	if err := writeStatus(&b, "PEL", "ready"); err != nil || b.String() != "PEL: ready\n" {
		t.Errorf("writeStatus(PEL, ready) = %q, %v, want PEL: ready newline, nil", b.String(), err)
	}
	problem := errors.New("report full")
	err := writeStatus(writerFunc(func(p []byte) (int, error) { return 0, problem }), "SQT", "ready")
	if !errors.Is(err, problem) {
		t.Errorf("writeStatus(failing writer) error = %v, want original error", err)
	}
}

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the supplied test writer.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

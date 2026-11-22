package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestCaptureText(t *testing.T) {
	var log bytes.Buffer
	got, err := captureText(strings.NewReader("recorded"), &log)
	if got != "recorded" || log.String() != got || err != nil {
		t.Errorf("captureText() = %q, %q, %v, want recorded twice, nil", got, log.String(), err)
	}
	problem := errors.New("recording stopped")
	log.Reset()
	got, err = captureText(readerFunc(func(p []byte) (int, error) { return copy(p, "partial"), problem }), &log)
	if got != "partial" || log.String() != got || !errors.Is(err, problem) {
		t.Errorf("captureText(failed source) = %q, %q, %v, want partial twice and original error", got, log.String(), err)
	}
	_, err = captureText(strings.NewReader("recorded"), writerFunc(func(p []byte) (int, error) { return 0, problem }))
	if !errors.Is(err, problem) {
		t.Errorf("captureText(failed log) error = %v, want original error", err)
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

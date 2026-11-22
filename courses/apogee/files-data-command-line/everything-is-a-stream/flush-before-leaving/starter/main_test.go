package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestBufferedReport(t *testing.T) {
	var out bytes.Buffer
	err := bufferedReport(&out, []string{"PEL ready", "SQT ready"})
	if err != nil || out.String() != "PEL ready\nSQT ready\n" {
		t.Errorf("bufferedReport() = %q, %v, want both lines and nil", out.String(), err)
	}
	problem := errors.New("destination stopped")
	for _, line := range []string{"small", strings.Repeat("x", 8192)} {
		err := bufferedReport(writerFunc(func(p []byte) (int, error) { return 0, problem }), []string{line})
		if !errors.Is(err, problem) {
			t.Errorf("bufferedReport(failing writer) error = %v, want original error", err)
		}
	}
}

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the supplied test writer.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

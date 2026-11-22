package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestMirrorReport(t *testing.T) {
	var screen, archive bytes.Buffer
	n, err := mirrorReport(strings.NewReader("PEL ready\n"), &screen, &archive)
	if n != 10 || screen.String() != "PEL ready\n" || archive.String() != screen.String() || err != nil {
		t.Errorf("mirrorReport() = %d, %q, %q, %v, want 10, identical report, nil", n, screen.String(), archive.String(), err)
	}
	problem := errors.New("screen stopped")
	archive.Reset()
	_, err = mirrorReport(strings.NewReader("report"), writerFunc(func(p []byte) (int, error) { return 0, problem }), &archive)
	if !errors.Is(err, problem) || archive.Len() != 0 {
		t.Errorf("mirrorReport(failed screen) error = %v, archive bytes %d, want original error and 0", err, archive.Len())
	}
	screen.Reset()
	_, err = mirrorReport(strings.NewReader("report"), &screen, writerFunc(func(p []byte) (int, error) { return 0, problem }))
	if !errors.Is(err, problem) || screen.String() != "report" {
		t.Errorf("mirrorReport(failed archive) screen = %q, error %v, want report and original error", screen.String(), err)
	}

}

// writerFunc adapts a test function to io.Writer.
type writerFunc func([]byte) (int, error)

// Write calls the supplied test writer.
func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

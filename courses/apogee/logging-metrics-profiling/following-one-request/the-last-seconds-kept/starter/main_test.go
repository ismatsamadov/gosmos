package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime/trace"
	"testing"
	"time"
)

// recorder starts a flight recorder that keeps a second of events, and stops it when the test ends.
func recorder(t *testing.T) *trace.FlightRecorder {
	t.Helper()
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{MinAge: time.Second})
	if err := fr.Start(); err != nil {
		t.Fatalf("starting the flight recorder: %v", err)
	}
	t.Cleanup(fr.Stop)
	return fr
}

// fly stands in for a pass: three goroutines hand the radio's lines along, then the pass ends with err.
func fly(err error) func() error {
	return func() error {
		lines := make(chan string)
		for _, l := range []string{"PING", "PING", "PING"} {
			go func() { lines <- l }()
		}
		for range 3 {
			<-lines
		}
		return err
	}
}

var noAck = errors.New("pass cirrine-1/LAN/17759: no ack before LOS")

// TestFailedPass keeps flight.trace in the lesson folder: open it with go tool trace afterwards.
func TestFailedPass(t *testing.T) {
	fr := recorder(t)
	os.Remove("flight.trace")
	err := watchPass(fr, "flight.trace", fly(noAck))
	if !errors.Is(err, noAck) {
		t.Errorf(`watchPass(fr, "flight.trace", a failed pass) = %v, want %q, the pass's error`, err, noAck)
	}
	b, readErr := os.ReadFile("flight.trace")
	if readErr != nil {
		t.Fatalf(`watchPass(fr, "flight.trace", a failed pass) left no flight.trace: %v`, readErr)
	}
	if head := b[:min(5, len(b))]; !bytes.Equal(head, []byte("go 1.")) {
		t.Errorf(`flight.trace starts with %q, want "go 1.", the start of every trace`, head)
	}
	if !fr.Enabled() {
		t.Errorf(`fr.Enabled() after watchPass(fr, "flight.trace", a failed pass) = false, want true: the recorder keeps recording`)
	}
}

func TestFlownPass(t *testing.T) {
	fr := recorder(t)
	path := filepath.Join(t.TempDir(), "flight.trace")
	if err := watchPass(fr, path, fly(nil)); err != nil {
		t.Errorf("watchPass(fr, path, a flown pass) = %v, want nil", err)
	}
	if _, err := os.Stat(path); err == nil {
		t.Errorf("watchPass(fr, path, a flown pass) wrote %s, want no file: only a failure is kept", filepath.Base(path))
	}
	if !fr.Enabled() {
		t.Errorf("fr.Enabled() after watchPass(fr, path, a flown pass) = false, want true")
	}
}

func TestMissingFolder(t *testing.T) {
	fr := recorder(t)
	path := filepath.Join(t.TempDir(), "no such folder", "flight.trace")
	err := watchPass(fr, path, fly(noAck))
	if !errors.Is(err, noAck) {
		t.Errorf("watchPass(fr, a path in a missing folder, a failed pass) = %v, want the pass's error first", err)
	}
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Errorf("watchPass(fr, a path in a missing folder, a failed pass) = %v, want the error from creating the file too", err)
	}
}

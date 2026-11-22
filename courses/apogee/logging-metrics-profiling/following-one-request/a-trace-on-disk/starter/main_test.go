package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
	"testing"
)

// replay stands in for the bench replay: the night's three passes, each on a goroutine of its
// own, handing its pass ID back when it is done.
func replay() {
	done := make(chan string)
	for _, id := range []string{"cirrine-1/LAN/17759", "furrow-b/LAN/22379", "furrow-a/LAN/22382"} {
		go func() { done <- id }()
	}
	for range 3 {
		<-done
	}
}

// stopped reports whether no trace is running, stopping one that is so later tests can start theirs.
func stopped() bool {
	if err := trace.Start(io.Discard); err != nil {
		trace.Stop()
		return false
	}
	trace.Stop()
	return true
}

// TestTraceTo writes trace.out in the lesson folder: open it with go tool trace afterwards.
func TestTraceTo(t *testing.T) {
	os.Remove("trace.out")
	runs := 0
	if err := traceTo("trace.out", func() { runs++; replay() }); err != nil {
		t.Fatalf(`traceTo("trace.out", replay) = %v, want nil`, err)
	}
	if runs != 1 {
		t.Errorf(`traceTo("trace.out", replay) ran the replay %d times, want 1`, runs)
	}
	if !stopped() {
		t.Errorf(`trace.Start after traceTo("trace.out", replay) failed, want it to work: the trace was never stopped`)
	}
	b, err := os.ReadFile("trace.out")
	if err != nil {
		t.Fatalf(`traceTo("trace.out", replay) left no file to read: %v`, err)
	}
	head := b[:min(16, len(b))]
	if !bytes.HasPrefix(head, []byte("go 1.")) || !bytes.Contains(head, []byte(" trace")) {
		t.Errorf(`trace.out starts with %q, want "go 1." and " trace", the start of every trace`, head)
	}
}

func TestTraceToMissingFolder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no such folder", "trace.out")
	runs := 0
	err := traceTo(path, func() { runs++ })
	if err == nil {
		t.Errorf("traceTo(a path in a missing folder, replay) = nil, want the error from creating the file")
	}
	if runs != 0 {
		t.Errorf("traceTo(a path in a missing folder, replay) ran the replay %d times, want 0: nothing starts without a file", runs)
	}
	if !stopped() {
		t.Errorf("trace.Start after traceTo(a path in a missing folder, replay) failed, want it to work: a trace was left running")
	}
}

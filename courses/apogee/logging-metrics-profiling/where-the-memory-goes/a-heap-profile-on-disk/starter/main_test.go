package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// collections returns how many garbage collections the program has finished so far.
func collections() uint32 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.NumGC
}

// TestHeapProfile saves heap.prof in the lesson folder, as the runbook's first line saves
// uplinkd's: you can open it with go tool pprof afterwards.
func TestHeapProfile(t *testing.T) {
	os.Remove("heap.prof")
	before := collections()
	if err := heapProfile("heap.prof"); err != nil {
		t.Fatalf(`heapProfile("heap.prof") = %v, want nil`, err)
	}
	ran := collections() - before
	b, err := os.ReadFile("heap.prof")
	if err != nil {
		t.Fatalf(`heapProfile("heap.prof") left no file to read: %v`, err)
	}
	switch {
	case len(b) == 0:
		t.Errorf(`heapProfile("heap.prof") wrote an empty file, want a heap profile`)
	case !bytes.HasPrefix(b, []byte{0x1f, 0x8b}):
		t.Errorf(`heapProfile("heap.prof") wrote a file starting % x, want 1f 8b, the first bytes of a profile written with WriteTo(f, 0)`, b[:min(2, len(b))])
	}
	if ran == 0 {
		t.Errorf(`heapProfile("heap.prof") ran 0 garbage collections, want 1 before writing: the profile counts the heap as of the last one`)
	}
}

func TestHeapProfileMissingFolder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no such folder", "heap.prof")
	if err := heapProfile(path); err == nil {
		t.Errorf("heapProfile(a path in a missing folder) = nil, want the error from creating the file")
	}
}

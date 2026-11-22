package main

import (
	"os"
	"runtime"
	"runtime/pprof"
)

// heapProfile writes the heap in use, as of a fresh garbage collection, to a new file at path.
func heapProfile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, _ = runtime.GC, pprof.Lookup // TODO: collect, then write the heap profile into f
	return f.Close()
}

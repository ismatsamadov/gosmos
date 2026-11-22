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
	runtime.GC()
	if err := pprof.Lookup("heap").WriteTo(f, 0); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

package main

import (
	"os"
	"runtime/trace"
)

// traceTo runs work while an execution trace of it is written to a new file at path.
func traceTo(path string, work func()) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := trace.Start(f); err != nil {
		f.Close()
		return err
	}
	work()
	trace.Stop()
	return f.Close()
}

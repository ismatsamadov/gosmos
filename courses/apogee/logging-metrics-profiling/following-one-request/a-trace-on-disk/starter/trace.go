package main

import (
	"os"
	"runtime/trace"
)

// traceTo runs work while an execution trace of it is written to a new file at path.
func traceTo(path string, work func()) error {
	_, _ = os.Create, trace.Start // TODO: create the file, trace work into it, then stop and close
	work()
	return nil
}

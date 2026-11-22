package main

import (
	"os"
	"runtime/pprof"
)

// profileTo writes a CPU profile of work to the file at path.
func profileTo(path string, work func()) error {
	_, _ = os.Create, pprof.StartCPUProfile // TODO: profile work into the file at path
	work()
	return nil
}

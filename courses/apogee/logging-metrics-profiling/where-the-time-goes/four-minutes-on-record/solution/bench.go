package main

import (
	"os"
	"runtime/pprof"
)

// profileTo writes a CPU profile of work to the file at path.
func profileTo(path string, work func()) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return err
	}
	work()
	pprof.StopCPUProfile()
	return f.Close()
}

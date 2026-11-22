package main

import (
	"os"
	"runtime/trace"
)

// watchPass runs one pass and, when it fails, saves the flight recorder's last seconds to path.
func watchPass(fr *trace.FlightRecorder, path string, run func() error) error {
	_ = os.Create // TODO: when run fails, write fr's snapshot into a new file at path
	return run()
}

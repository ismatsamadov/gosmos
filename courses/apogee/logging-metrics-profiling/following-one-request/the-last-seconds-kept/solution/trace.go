package main

import (
	"errors"
	"os"
	"runtime/trace"
)

// watchPass runs one pass and, when it fails, saves the flight recorder's last seconds to path.
func watchPass(fr *trace.FlightRecorder, path string, run func() error) error {
	err := run()
	if err == nil {
		return nil
	}
	f, ferr := os.Create(path)
	if ferr != nil {
		return errors.Join(err, ferr)
	}
	_, werr := fr.WriteTo(f)
	return errors.Join(err, werr, f.Close())
}

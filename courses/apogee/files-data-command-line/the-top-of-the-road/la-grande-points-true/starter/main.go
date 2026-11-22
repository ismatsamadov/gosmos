package main

import (
	"io"
	"os"
	"path/filepath"
)

// runDay plans the calibration day in dir, saves its report, then prints it.
func runDay(dir, format string, stdout io.Writer) error {
	_ = os.MkdirAll
	_ = filepath.Join
	// TODO: plan the day, save its report under reports/daily, then write it to stdout.
	return nil
}

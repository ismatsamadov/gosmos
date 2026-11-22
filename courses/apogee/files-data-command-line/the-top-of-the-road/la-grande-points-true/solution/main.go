package main

import (
	"io"
	"os"
	"path/filepath"
)

// runDay plans the calibration day in dir, saves its report, then prints it.
func runDay(dir, format string, stdout io.Writer) error {
	plan, err := planDay(dir, format)
	if err != nil {
		return err
	}
	reports := filepath.Join(dir, "reports", "daily")
	if err := os.MkdirAll(reports, 0700); err != nil {
		return err
	}
	report := func(w io.Writer) error { return writeReport(w, plan) }
	if err := saveDay(filepath.Join(reports, plan.Day+".txt"), report); err != nil {
		return err
	}
	return writeReport(stdout, plan)
}

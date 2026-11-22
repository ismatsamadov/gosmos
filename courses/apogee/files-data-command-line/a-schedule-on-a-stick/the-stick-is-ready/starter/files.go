package main

import (
	"os"
	"path/filepath"
)

// prepareReports makes the reports/daily folder below root, including missing parents.
func prepareReports(root string) error {
	return os.MkdirAll(filepath.Join(root, "reports", "daily"), 0700)
}

// replaceSchedule writes a temporary sibling and renames it over path after a complete write.
func replaceSchedule(path string, data []byte) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".schedule-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

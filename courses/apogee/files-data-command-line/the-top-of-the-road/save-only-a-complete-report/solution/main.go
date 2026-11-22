package main

import (
	"io"
	"os"
	"path/filepath"
)

// saveDay writes the report that report produces to path, through a temporary sibling.
func saveDay(path string, report func(io.Writer) error) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".day-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if err := report(file); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

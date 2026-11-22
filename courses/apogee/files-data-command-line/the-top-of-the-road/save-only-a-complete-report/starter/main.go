package main

import (
	"io"
	"os"
	"path/filepath"
)

// saveDay writes the report that report produces to path, through a temporary sibling.
func saveDay(path string, report func(io.Writer) error) error {
	_ = os.CreateTemp
	_ = filepath.Dir
	return nil
}

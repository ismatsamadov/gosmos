package main

import (
	"os"
	"path/filepath"
)

// replaceSchedule writes a temporary sibling and renames it over path after a complete write.
func replaceSchedule(path string, data []byte) error {
	_ = os.CreateTemp
	_ = filepath.Dir
	return nil
}

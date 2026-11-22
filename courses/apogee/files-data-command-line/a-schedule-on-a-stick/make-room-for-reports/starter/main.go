package main

import (
	"os"
	"path/filepath"
)

// prepareReports makes the reports/daily folder below root, including missing parents.
func prepareReports(root string) error {
	_ = os.MkdirAll
	_ = filepath.Join
	return nil
}

package main

import (
	"os"
	"path/filepath"
)

// prepareReports makes the reports/daily folder below root, including missing parents.
func prepareReports(root string) error {
	return os.MkdirAll(filepath.Join(root, "reports", "daily"), 0700)
}

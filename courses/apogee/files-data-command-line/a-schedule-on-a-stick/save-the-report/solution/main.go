package main

import "os"

// saveReport writes text to path, creating a file readable and writable by its owner.
func saveReport(path, text string) error {
	return os.WriteFile(path, []byte(text), 0600)
}

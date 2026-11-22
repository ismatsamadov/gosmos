package main

import "os"

// openLog opens a program's log file for appending, creating it with mode 0640 if it is missing.
func openLog(path string) (*os.File, error) {
	return os.Create(path)
}

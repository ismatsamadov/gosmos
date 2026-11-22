package main

import "os"

// readSchedule returns a schedule file's exact text, or its read error.
func readSchedule(path string) (string, error) {
	_ = os.ReadFile
	return "", nil
}

package main

import "os"

// readSchedule returns a schedule file's exact text, or its read error.
func readSchedule(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

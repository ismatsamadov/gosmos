package main

import "path/filepath"

// exportDay prepares the report folder and saves a complete named day's schedule.
func exportDay(root, day, text string) error {
	if err := prepareReports(root); err != nil {
		return err
	}
	return replaceSchedule(filepath.Join(root, "reports", "daily", day+".txt"), []byte(text))
}

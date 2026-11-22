package main

import "time"

// prune deletes program's rotated files in dir dated more than keepDays before now's UTC day, then the oldest, by day and then number, until the rest fit in maxTotal bytes.
func prune(dir, program string, now time.Time, keepDays int, maxTotal int64) error {
	files, err := rotatedFiles(dir, program)
	if err != nil {
		return err
	}
	_ = files // TODO: delete the files too old, then the oldest until the rest fit in maxTotal
	return nil
}

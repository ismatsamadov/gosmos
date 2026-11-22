package main

import (
	"cmp"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// prune deletes program's rotated files in dir dated more than keepDays before now's UTC day, then the oldest, by day and then number, until the rest fit in maxTotal bytes.
func prune(dir, program string, now time.Time, keepDays int, maxTotal int64) error {
	files, err := rotatedFiles(dir, program)
	if err != nil {
		return err
	}
	slices.SortFunc(files, func(a, b rotated) int {
		return cmp.Or(strings.Compare(a.Day, b.Day), cmp.Compare(a.N, b.N))
	})
	var total int64
	for _, f := range files {
		total += f.Size
	}
	oldest := now.UTC().AddDate(0, 0, -keepDays).Format(time.DateOnly)
	for _, f := range files {
		if f.Day >= oldest && total <= maxTotal {
			break
		}
		if err := os.Remove(filepath.Join(dir, f.Name)); err != nil {
			return err
		}
		total -= f.Size
	}
	return nil
}

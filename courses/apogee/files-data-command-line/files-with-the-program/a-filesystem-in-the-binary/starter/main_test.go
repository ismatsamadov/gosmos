package main

import (
	"errors"
	"io/fs"
	"testing"
)

func TestLoadTemplate(t *testing.T) {
	// An empty folder: templates read from the disk at run time are not here.
	t.Chdir(t.TempDir())
	for _, tc := range []struct{ name, file, want string }{
		{"day report", "day.txt", "Calibration day {day} at {station}\n"},
		{"pass line", "pass.txt", "{aos}-{los} {sat} {id}\n"},
		{"nested stow notice", "notices/stow.txt", "{station} stowed: {reason}\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadTemplate(tc.file)
			if got != tc.want || err != nil {
				t.Errorf("loadTemplate(%q) = %q, %v, want %q, nil", tc.file, got, err, tc.want)
			}
		})
	}
	got, err := loadTemplate("night.txt")
	if got != "" || !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("loadTemplate(%q) = %q, %v, want \"\" and a not-exist error", "night.txt", got, err)
	}
}

package main

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestTemplateFS(t *testing.T) {
	view, err := templateFS(resources)
	if view == nil || err != nil {
		t.Fatalf("templateFS(resources) = %v, %v, want a view and nil", view, err)
	}
	for _, tc := range []struct {
		name, file string
		fields     map[string]string
		want       string
	}{
		{"day report", "day.txt", map[string]string{"day": "2027-10-14", "station": "SQT"}, "Calibration day 2027-10-14 at SQT\n"},
		{"nested stow notice", "notices/stow.txt", map[string]string{"station": "CWR", "reason": "wind 81 km/h"}, "CWR stowed: wind 81 km/h\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := render(view, tc.file, tc.fields)
			if got != tc.want || err != nil {
				t.Errorf("render(view, %q) = %q, %v, want %q, nil", tc.file, got, err, tc.want)
			}
		})
	}
	if _, err := fs.ReadFile(view, "templates/day.txt"); err == nil {
		t.Errorf("fs.ReadFile(view, %q) error = nil, want an error: the view starts inside templates", "templates/day.txt")
	}
	if _, err := render(view, "night.txt", nil); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("render(view, %q) error = %v, want a not-exist error", "night.txt", err)
	}

	_, err = templateFS(fstest.MapFS{"stations.csv": {Data: []byte("code\n")}})
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("templateFS(no templates folder) error = %v, want a not-exist error", err)
	}
}

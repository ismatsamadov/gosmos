package main

import (
	"errors"
	"io/fs"
	"slices"
	"testing"
	"testing/fstest"
)

func TestListTemplates(t *testing.T) {
	want := []string{"templates/day.txt", "templates/notices/handover.txt", "templates/notices/stow.txt", "templates/pass.txt"}
	got, err := listTemplates(resources)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("listTemplates(resources) = %q, %v, want %q, nil", got, err, want)
	}

	drafts := fstest.MapFS{
		"templates/wind/cwr/gusts.txt": {Data: []byte("Gusts {kmh} km/h\n")},
		"templates/drafts":             {Mode: fs.ModeDir}, // an empty folder
		"templates/day.txt":            {Data: []byte("Day {day}\n")},
		"stations.csv":                 {Data: []byte("code,name,lat_deg,lon_deg,mask_deg,bands\n")},
	}
	want = []string{"templates/day.txt", "templates/wind/cwr/gusts.txt"}
	got, err = listTemplates(drafts)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("listTemplates(deeper folders and an empty one) = %q, %v, want %q, nil", got, err, want)
	}

	_, err = listTemplates(fstest.MapFS{"stations.csv": {Data: []byte("code\n")}})
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("listTemplates(no templates folder) error = %v, want a not-exist error", err)
	}
}

package main

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

// stick returns a stick holding calibration.json with the given text.
func stick(text string) fstest.MapFS {
	return fstest.MapFS{"calibration.json": {Data: []byte(text)}}
}

// deniedFS is a stick the laptop is not allowed to read.
type deniedFS struct{}

// Open refuses every name.
func (deniedFS) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
}

func TestLoadResources(t *testing.T) {
	for _, tc := range []struct {
		name     string
		override fs.FS
		want     Offsets
	}{
		{"empty stick", fstest.MapFS{}, Offsets{}},
		{"stick with notes only", fstest.MapFS{"notes.txt": {Data: []byte("gloves\n")}}, Offsets{}},
		{"port ellery check", stick(`{"az_offset_deg":0.04,"el_offset_deg":-0.02}`), Offsets{AzOffsetDeg: 0.04, ElOffsetDeg: -0.02}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadResources(resources, tc.override)
			if len(got.Stations) != 5 || got.Reference != tc.want || err != nil {
				t.Errorf("loadResources(resources, %s) = %d stations, %+v, %v, want 5 stations, %+v, nil", tc.name, len(got.Stations), got.Reference, err, tc.want)
			}
		})
	}
}

func TestLoadResourcesSurfacesErrors(t *testing.T) {
	table, err := resources.ReadFile("stations.csv")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name           string
		fsys, override fs.FS
		target         error // nil when any error will do
	}{
		{"file cut off on the stick", resources, stick(`{"az_offset_deg":0.04,`), nil},
		{"misspelled field on the stick", resources, stick(`{"az_offset":0.04,"el_offset_deg":-0.02}`), nil},
		{"stick the laptop may not read", resources, deniedFS{}, fs.ErrPermission},
		{"no default either", fstest.MapFS{"stations.csv": {Data: table}}, fstest.MapFS{}, fs.ErrNotExist},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadResources(tc.fsys, tc.override)
			if err == nil || tc.target != nil && !errors.Is(err, tc.target) {
				t.Errorf("loadResources(%s) = %+v, %v, want an error, not a fallback to zero", tc.name, got.Reference, err)
			}
		})
	}
}

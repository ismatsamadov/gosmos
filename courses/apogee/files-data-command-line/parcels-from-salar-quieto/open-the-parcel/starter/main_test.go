package main

import (
	"archive/zip"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

const (
	// manifest is manifest.xml as La Grande's dish controller packed it.
	manifest = `<?xml version="1.0" encoding="UTF-8"?>
<parcel version="1">
  <dish>LG</dish>
  <firmware>5.3</firmware>
  <made>2027-10-15T09:40:00Z</made>
  <files>
    <file size="153">correction.xml</file>
    <file size="136">log/2027-10-14.txt</file>
  </files>
</parcel>
`
	// correction is La Grande's correction, as the parcel holds it.
	correction = `<?xml version="1.0" encoding="UTF-8"?>
<correction dish="LG">
  <az_offset_deg>0.18</az_offset_deg>
  <el_offset_deg>-0.06</el_offset_deg>
</correction>
`
	// day is the last day of the controller log, as the parcel holds it.
	day = "# night shift 2027-10-14\n" +
		"2027-10-14T09:12:04Z\tLG\t201.40\t12.85\n" +
		"2027-10-14T09:12:05Z\tLG\t201.62\t13.10\n" +
		"2027-10-14T09:12:06Z\tLG\t201.85\t13.36\n"
)

// file is one file to pack into a test parcel.
type file struct {
	name, text string
}

// writeParcel packs files, in order, into a zip in a new folder and returns the zip's path.
func writeParcel(t *testing.T, files ...file) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "lg-settings-2027-10-15.zip")
	out, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(out)
	for _, f := range files {
		w, err := zw.Create(f.name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(f.text)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

// same reports whether two manifests hold the same values.
func same(a, b Manifest) bool {
	return a.Version == b.Version && a.Dish == b.Dish && a.Made.Equal(b.Made) && slices.Equal(a.Files, b.Files)
}

func TestParcelManifest(t *testing.T) {
	want := Manifest{
		Version: 1,
		Dish:    "LG",
		Made:    time.Date(2027, 10, 15, 9, 40, 0, 0, time.UTC),
		Files: []Entry{
			{Name: "correction.xml", Size: 153},
			{Name: "log/2027-10-14.txt", Size: 136},
		},
	}
	for _, tc := range []struct {
		name  string
		files []file
	}{
		{"the controller's parcel", []file{{"manifest.xml", manifest}, {"correction.xml", correction}, {"log/2027-10-14.txt", day}}},
		{"a parcel with its manifest last", []file{{"correction.xml", correction}, {"log/2027-10-14.txt", day}, {"manifest.xml", manifest}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parcelManifest(writeParcel(t, tc.files...))
			if err != nil || !same(got, want) {
				t.Errorf("parcelManifest(%s) = %+v, %v, want %+v, nil", tc.name, got, err, want)
			}
		})
	}
}

func TestParcelManifestRefuses(t *testing.T) {
	notZip := filepath.Join(t.TempDir(), "lg-settings-2027-10-15.zip")
	if err := os.WriteFile(notZip, []byte(day), 0o600); err != nil {
		t.Fatal(err)
	}
	noDish := strings.Replace(manifest, "  <dish>LG</dish>\n", "", 1)
	for _, tc := range []struct {
		name, path string
	}{
		{"a parcel with no manifest.xml", writeParcel(t, file{"correction.xml", correction}, file{"log/2027-10-14.txt", day})},
		{"a parcel whose manifest has no <dish>", writeParcel(t, file{"manifest.xml", noDish}, file{"correction.xml", correction})},
		{"a log file named like a parcel", notZip},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parcelManifest(tc.path)
			if err == nil {
				t.Errorf("parcelManifest(%s) = %+v, nil, want an error", tc.name, got)
			}
		})
	}
	missing := filepath.Join(t.TempDir(), "lg-settings-2027-10-16.zip")
	if _, err := parcelManifest(missing); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("parcelManifest(a parcel not on the stick) error = %v, want a not-exist error", err)
	}
}

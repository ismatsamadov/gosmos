package main

import (
	"archive/zip"
	"bytes"
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
	name string
	data []byte
}

// settings is the parcel of La Grande's settings, as the controller packed it.
var settings = []file{
	{"manifest.xml", []byte(manifest)},
	{"correction.xml", []byte(correction)},
	{"log/2027-10-14.txt", []byte(day)},
}

// writeParcel packs files, in order, into a zip in a folder of its own and returns the zip's path.
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
		if _, err := w.Write(f.data); err != nil {
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

// stick makes a test folder holding stick, the folder to unpack into, and two empty folders beside it:
// outside, and home, which becomes the current folder. It returns the three.
func stick(t *testing.T) (dir, outside, home string) {
	t.Helper()
	base := t.TempDir()
	dir = filepath.Join(base, "stick")
	outside = filepath.Join(base, "outside")
	home = filepath.Join(base, "home")
	for _, d := range []string{dir, outside, home} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(home)
	return dir, outside, home
}

// tree returns the path of every file and link below dir, relative to it, in walk order.
func tree(t *testing.T, dir string) []string {
	t.Helper()
	var paths []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, err := filepath.Rel(dir, p)
			if err != nil {
				return err
			}
			paths = append(paths, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return paths
}

// same reports whether two manifests hold the same values.
func same(a, b Manifest) bool {
	return a.Version == b.Version && a.Dish == b.Dish && a.Made.Equal(b.Made) && slices.Equal(a.Files, b.Files)
}

func TestUnpackParcel(t *testing.T) {
	dir, _, home := stick(t)
	want := Manifest{
		Version: 1,
		Dish:    "LG",
		Made:    time.Date(2027, 10, 15, 9, 40, 0, 0, time.UTC),
		Files: []Entry{
			{Name: "correction.xml", Size: 153},
			{Name: "log/2027-10-14.txt", Size: 136},
		},
	}
	got, err := unpackParcel(writeParcel(t, settings...), dir)
	if err != nil || !same(got, want) {
		t.Fatalf("unpackParcel(the settings parcel) = %+v, %v, want %+v, nil", got, err, want)
	}
	for _, f := range settings {
		data, err := os.ReadFile(filepath.Join(dir, f.name))
		if err != nil {
			t.Errorf("after unpackParcel(the settings parcel), the stick has no %s, want the file there", f.name)
		} else if !bytes.Equal(data, f.data) {
			t.Errorf("after unpackParcel(the settings parcel), %s holds %q, want %q", f.name, data, f.data)
		}
	}
	if got := tree(t, home); len(got) > 0 {
		t.Errorf("the current folder gained %q, want every file below the stick's folder", got)
	}
}

func TestUnpackParcelWritesNothingUnsafe(t *testing.T) {
	dir, outside, home := stick(t)
	path := writeParcel(t, append(slices.Clone(settings), file{"../gctl", []byte("#!/bin/sh\n")})...)
	if _, err := unpackParcel(path, dir); err == nil || !strings.Contains(err.Error(), "../gctl") {
		t.Errorf("unpackParcel(a parcel holding ../gctl) error = %v, want one naming ../gctl", err)
	}
	if got := tree(t, dir); len(got) > 0 {
		t.Errorf("unpackParcel(a parcel holding ../gctl) wrote %q on the stick, want nothing: check every name first", got)
	}
	want := []string{"home", "outside", "stick"}
	entries, err := os.ReadDir(filepath.Dir(dir))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	if !slices.Equal(got, want) {
		t.Errorf("unpackParcel(a parcel holding ../gctl) left the stick's parent folder holding %q, want %q", got, want)
	}
	if got := append(tree(t, outside), tree(t, home)...); len(got) > 0 {
		t.Errorf("unpackParcel(a parcel holding ../gctl) wrote %q outside the stick, want nothing", got)
	}
}

func TestUnpackParcelStaysInside(t *testing.T) {
	dir, outside, home := stick(t)
	if err := os.Symlink(filepath.Join("..", "outside"), filepath.Join(dir, "log")); err != nil {
		t.Fatal(err)
	}
	if _, err := unpackParcel(writeParcel(t, settings...), dir); err == nil {
		t.Errorf("unpackParcel(the settings parcel, with log a link out of the stick) error = nil, want an error")
	}
	if got := append(tree(t, outside), tree(t, home)...); len(got) > 0 {
		t.Errorf("unpackParcel(the settings parcel, with log a link out of the stick) wrote %q outside the stick, want nothing", got)
	}
}

func TestUnpackParcelCapsEveryFile(t *testing.T) {
	dir, _, _ := stick(t)
	swollen := file{"log/2027-10-15.txt", make([]byte, 16<<20)} // 16 MiB of zeros
	if _, err := unpackParcel(writeParcel(t, append(slices.Clone(settings), swollen)...), dir); err == nil {
		t.Errorf("unpackParcel(a parcel with 16 MiB of zeros) error = nil, want an error")
	}
	for _, name := range tree(t, dir) {
		info, err := os.Stat(filepath.Join(dir, name))
		if err == nil && info.Size() > maxFileBytes+1 {
			t.Errorf("unpackParcel(a parcel with 16 MiB of zeros) left %s at %d bytes, want at most one byte past maxFileBytes", name, info.Size())
		}
	}
}

func TestUnpackParcelRefuses(t *testing.T) {
	notZip := filepath.Join(t.TempDir(), "lg-settings-2027-10-15.zip")
	if err := os.WriteFile(notZip, []byte(day), 0o600); err != nil {
		t.Fatal(err)
	}
	noDish := strings.Replace(manifest, "  <dish>LG</dish>\n", "", 1)
	for _, tc := range []struct {
		name, path string
	}{
		{"a parcel with no manifest.xml", writeParcel(t, settings[1:]...)},
		{"a parcel whose manifest has no <dish>", writeParcel(t, file{"manifest.xml", []byte(noDish)}, settings[1])},
		{"a log file named like a parcel", notZip},
		{"a parcel that is not on the stick", filepath.Join(t.TempDir(), "lg-settings-2027-10-16.zip")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, _, _ := stick(t)
			got, err := unpackParcel(tc.path, dir)
			if err == nil {
				t.Errorf("unpackParcel(%s) = %+v, nil, want an error", tc.name, got)
			}
		})
	}
}

package main

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const (
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

// zipped packs files, in order, into a zip in memory and returns the zip's list of files.
func zipped(t *testing.T, files ...file) []*zip.File {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
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
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return zr.File
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

// openRoot opens dir as a root, closed when the test ends.
func openRoot(t *testing.T, dir string) *os.Root {
	t.Helper()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { root.Close() })
	return root
}

// names returns the names a folder holds.
func names(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range entries {
		out = append(out, e.Name())
	}
	return out
}

func TestUnpackFile(t *testing.T) {
	dir, _, home := stick(t)
	root := openRoot(t, dir)
	files := []file{{"correction.xml", []byte(correction)}, {"log/2027-10-14.txt", []byte(day)}}
	for i, f := range zipped(t, files...) {
		if err := unpackFile(root, f, 64<<10); err != nil {
			t.Fatalf("unpackFile(%s, 64 KiB) = %v, want nil", f.Name, err)
		}
		got, err := os.ReadFile(filepath.Join(dir, files[i].name))
		if err != nil {
			t.Errorf("unpackFile(%s, 64 KiB) = nil, but the stick has no %s, want the file there", f.Name, files[i].name)
		} else if !bytes.Equal(got, files[i].data) {
			t.Errorf("unpackFile(%s, 64 KiB) wrote %q, want %q", f.Name, got, files[i].data)
		}
	}
	if got := names(t, home); len(got) > 0 {
		t.Errorf("the current folder gained %q, want every file below the stick's folder", got)
	}
}

func TestUnpackFileStaysInside(t *testing.T) {
	t.Run("a link that leads out", func(t *testing.T) {
		dir, outside, home := stick(t)
		if err := os.Symlink(filepath.Join("..", "outside"), filepath.Join(dir, "log")); err != nil {
			t.Fatal(err)
		}
		f := zipped(t, file{"log/2027-10-14.txt", []byte(day)})[0]
		if err := unpackFile(openRoot(t, dir), f, 64<<10); err == nil {
			t.Errorf("unpackFile(log/2027-10-14.txt, with log a link out of the stick) = nil, want an error")
		}
		if got := names(t, outside); len(got) > 0 {
			t.Errorf("unpackFile(log/2027-10-14.txt, with log a link out of the stick) wrote %q outside the stick, want nothing", got)
		}
		if got := names(t, home); len(got) > 0 {
			t.Errorf("the current folder gained %q, want nothing", got)
		}
	})
	t.Run("a name that climbs out", func(t *testing.T) {
		dir, _, home := stick(t)
		f := zipped(t, file{"../gctl", []byte("#!/bin/sh\n")})[0]
		if err := unpackFile(openRoot(t, dir), f, 64<<10); err == nil {
			t.Errorf("unpackFile(../gctl) = nil, want an error")
		}
		want := []string{"home", "outside", "stick"}
		if got := names(t, filepath.Dir(dir)); !slices.Equal(got, want) {
			t.Errorf("unpackFile(../gctl) left the stick's parent folder holding %q, want %q", got, want)
		}
		if got := names(t, home); len(got) > 0 {
			t.Errorf("the current folder gained %q, want nothing", got)
		}
	})
}

func TestUnpackFileLimit(t *testing.T) {
	dir, _, _ := stick(t)
	root := openRoot(t, dir)
	limit := int64(len(day))
	f := zipped(t, file{"log/2027-10-14.txt", []byte(day)})[0]
	if err := unpackFile(root, f, limit); err != nil {
		t.Errorf("unpackFile(a day of exactly %d bytes, %d) = %v, want nil", limit, limit, err)
	}
	swollen := zipped(t, file{"log/2027-10-15.txt", make([]byte, 16<<20)})[0] // 16 MiB of zeros
	if err := unpackFile(root, swollen, 64<<10); err == nil {
		t.Errorf("unpackFile(16 MiB of zeros, 64 KiB) = nil, want an error")
	}
	info, err := os.Stat(filepath.Join(dir, "log", "2027-10-15.txt"))
	if err == nil && info.Size() > 64<<10+1 {
		t.Errorf("unpackFile(16 MiB of zeros, 64 KiB) left %d bytes on the stick, want at most one byte past 64 KiB", info.Size())
	}
}

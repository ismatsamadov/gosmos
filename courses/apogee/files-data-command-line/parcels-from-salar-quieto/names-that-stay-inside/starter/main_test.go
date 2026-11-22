package main

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

// zipped packs one empty file per name into a zip in memory, in order, and returns the zip's list of files.
func zipped(t *testing.T, names ...string) []*zip.File {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, name := range names {
		if _, err := zw.Create(name); err != nil {
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

func TestCheckNamesAccepts(t *testing.T) {
	for _, tc := range []struct {
		name  string
		names []string
	}{
		{"the controller's parcel", []string{"manifest.xml", "correction.xml", "log/2027-10-14.txt"}},
		{"a name that wanders but stays inside", []string{"manifest.xml", "log/../correction.xml"}},
		{"two dots inside a file name", []string{"manifest.xml", "log/2027-10-14..txt"}},
		{"an empty parcel", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := checkNames(zipped(t, tc.names...)); err != nil {
				t.Errorf("checkNames(%q) = %v, want nil", tc.names, err)
			}
		})
	}
}

func TestCheckNamesRefuses(t *testing.T) {
	for _, tc := range []struct {
		name  string
		names []string
		bad   string
	}{
		{"an absolute name", []string{"manifest.xml", "/etc/hosts"}, "/etc/hosts"},
		{"a name that climbs out", []string{"manifest.xml", "correction.xml", "../gctl"}, "../gctl"},
		{"a name that wanders, then climbs out", []string{"manifest.xml", "log/../../.profile"}, "log/../../.profile"},
		{"two unsafe names", []string{"correction.xml", "../gctl", "/etc/hosts"}, "../gctl"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := checkNames(zipped(t, tc.names...))
			if err == nil || !strings.Contains(err.Error(), tc.bad) {
				t.Errorf("checkNames(%q) = %v, want an error naming %q", tc.names, err, tc.bad)
			}
		})
	}
}

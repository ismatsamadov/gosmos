package main

import (
	"slices"
	"strings"
	"testing"
	"time"
)

// manifest is manifest.xml as La Grande's dish controller packed it, <firmware> and all.
const manifest = `<?xml version="1.0" encoding="UTF-8"?>
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

// correction is the parcel's correction.xml, the wrong file to hand readManifest.
const correction = `<?xml version="1.0" encoding="UTF-8"?>
<correction dish="LG">
  <az_offset_deg>0.18</az_offset_deg>
  <el_offset_deg>-0.06</el_offset_deg>
</correction>
`

// same reports whether two manifests hold the same values.
func same(a, b Manifest) bool {
	return a.Version == b.Version && a.Dish == b.Dish && a.Made.Equal(b.Made) && slices.Equal(a.Files, b.Files)
}

func TestReadManifest(t *testing.T) {
	want := Manifest{
		Version: 1,
		Dish:    "LG",
		Made:    time.Date(2027, 10, 15, 9, 40, 0, 0, time.UTC),
		Files: []Entry{
			{Name: "correction.xml", Size: 153},
			{Name: "log/2027-10-14.txt", Size: 136},
		},
	}
	got, err := readManifest([]byte(manifest))
	if err != nil || !same(got, want) {
		t.Errorf("readManifest(manifest.xml with <firmware>) = %+v, %v, want %+v, nil", got, err, want)
	}
}

func TestReadManifestRefusesGaps(t *testing.T) {
	for _, tc := range []struct {
		name, data, missing string
	}{
		{"no <dish>", strings.Replace(manifest, "  <dish>LG</dish>\n", "", 1), "dish"},
		{"an empty <dish>", strings.Replace(manifest, "<dish>LG</dish>", "<dish></dish>", 1), "dish"},
		{"no <made>", strings.Replace(manifest, "  <made>2027-10-15T09:40:00Z</made>\n", "", 1), "made"},
		{"no <file>", `<parcel version="1"><dish>LG</dish><made>2027-10-15T09:40:00Z</made><files></files></parcel>`, "file"},
		{"correction.xml by mistake", correction, "dish"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readManifest([]byte(tc.data))
			if err == nil || !strings.Contains(err.Error(), tc.missing) {
				t.Errorf("readManifest(%s) error = %v, want an error naming %s", tc.name, err, tc.missing)
			}
		})
	}
}

func TestReadManifestRefusesBrokenXML(t *testing.T) {
	for _, tc := range []struct {
		name, data string
	}{
		{"cut off halfway", manifest[:len(manifest)/2]},
		{"a <made> that is not a time", strings.Replace(manifest, "2027-10-15T09:40:00Z", "Friday night", 1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readManifest([]byte(tc.data))
			if err == nil {
				t.Errorf("readManifest(%s) error = nil, want an error", tc.name)
			}
		})
	}
}

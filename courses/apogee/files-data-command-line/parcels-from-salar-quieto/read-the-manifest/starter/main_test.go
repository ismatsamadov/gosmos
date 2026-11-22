package main

import (
	"encoding/xml"
	"slices"
	"testing"
	"time"
)

// manifest is manifest.xml as La Grande's dish controller packed it into the parcel.
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

func TestManifestDecodes(t *testing.T) {
	var got Manifest
	if err := xml.Unmarshal([]byte(manifest), &got); err != nil {
		t.Fatalf("xml.Unmarshal(manifest.xml) error = %v, want nil", err)
	}
	if got.Version != 1 {
		t.Errorf("xml.Unmarshal(manifest.xml) gives Version = %d, want 1, the version attribute of <parcel>", got.Version)
	}
	if got.Dish != "LG" {
		t.Errorf("xml.Unmarshal(manifest.xml) gives Dish = %q, want \"LG\", the text of <dish>", got.Dish)
	}
	made := time.Date(2027, 10, 15, 9, 40, 0, 0, time.UTC)
	if !got.Made.Equal(made) {
		t.Errorf("xml.Unmarshal(manifest.xml) gives Made = %v, want %v, the text of <made>", got.Made, made)
	}
	files := []Entry{
		{Name: "correction.xml", Size: 153},
		{Name: "log/2027-10-14.txt", Size: 136},
	}
	if !slices.Equal(got.Files, files) {
		t.Errorf("xml.Unmarshal(manifest.xml) gives Files = %+v, want %+v, one per <file> inside <files>", got.Files, files)
	}
}

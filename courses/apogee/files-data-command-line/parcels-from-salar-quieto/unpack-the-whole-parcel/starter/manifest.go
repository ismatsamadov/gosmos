package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

// Manifest is a parcel's manifest.xml: the dish it is about, when the controller made it and the files it holds.
type Manifest struct {
	Version int       `xml:"version,attr"`
	Dish    string    `xml:"dish"`
	Made    time.Time `xml:"made"`
	Files   []Entry   `xml:"files>file"`
}

// Entry is one file a manifest lists: its name in the parcel and its size in bytes.
type Entry struct {
	Name string `xml:",chardata"`
	Size int64  `xml:"size,attr"`
}

// readManifest decodes manifest.xml and refuses one that lacks its dish, the time it was made or any file.
func readManifest(data []byte) (Manifest, error) {
	var m Manifest
	if err := xml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("manifest: %w", err)
	}
	if m.Dish == "" {
		return Manifest{}, errors.New("manifest: no <dish>")
	}
	if m.Made.IsZero() {
		return Manifest{}, errors.New("manifest: no <made>")
	}
	if len(m.Files) == 0 {
		return Manifest{}, errors.New("manifest: no <file>")
	}
	return m, nil
}

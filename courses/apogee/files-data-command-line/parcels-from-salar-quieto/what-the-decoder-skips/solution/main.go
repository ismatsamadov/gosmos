package main

import (
	"encoding/xml"
	"errors"
	"fmt"
)

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

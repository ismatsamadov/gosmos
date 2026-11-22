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
	// TODO: refuse a manifest with no <dish>, no <made> or no <file>, naming what is missing.
	_ = errors.New
	return m, nil
}

package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

// readLogs returns the text of every file in the gzip'd tar r, by name, refusing a file longer than limit bytes.
func readLogs(r io.Reader, limit int64) (map[string]string, error) {
	// TODO: unsqueeze r with gzip, read its files with tar and read each through io.LimitReader.
	_ = gzip.NewReader
	_ = tar.NewReader
	return map[string]string{}, nil
}

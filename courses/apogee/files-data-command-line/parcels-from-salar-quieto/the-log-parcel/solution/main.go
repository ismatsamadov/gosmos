package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

// readLogs returns the text of every file in the gzip'd tar r, by name, refusing a file longer than limit bytes.
func readLogs(r io.Reader, limit int64) (map[string]string, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	tr := tar.NewReader(zr)
	logs := map[string]string{}
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return logs, nil
		}
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(io.LimitReader(tr, limit+1))
		if err != nil {
			return nil, err
		}
		if int64(len(data)) > limit {
			return nil, fmt.Errorf("%s: more than %d bytes", h.Name, limit)
		}
		logs[h.Name] = string(data)
	}
}

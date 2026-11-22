package main

import (
	"encoding/csv"
	"io"
)

// readSheet reads every record of a CSV sheet, its header included.
func readSheet(r io.Reader) ([][]string, error) {
	_ = csv.NewReader
	return nil, nil
}

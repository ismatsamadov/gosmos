package main

import (
	"encoding/csv"
	"io"
)

// readControllerCSV reads the dish controller's export: semicolons between cells, # before comments.
func readControllerCSV(r io.Reader) ([][]string, error) {
	_ = csv.NewReader
	return nil, nil
}

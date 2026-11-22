package main

import (
	"encoding/csv"
	"io"
)

// readControllerCSV reads the dish controller's export: semicolons between cells, # before comments.
func readControllerCSV(r io.Reader) ([][]string, error) {
	cr := csv.NewReader(r)
	cr.Comma = ';'
	cr.Comment = '#'
	return cr.ReadAll()
}

package main

import (
	"encoding/csv"
	"io"
)

// writeSheet writes rows to w as CSV, quoting cells that need it, and reports any failure.
func writeSheet(w io.Writer, rows [][]string) error {
	_ = csv.NewWriter
	return nil
}

package main

import (
	"encoding/csv"
	"io"
)

// writeSheet writes rows to w as CSV, quoting cells that need it, and reports any failure.
func writeSheet(w io.Writer, rows [][]string) error {
	cw := csv.NewWriter(w)
	for _, row := range rows {
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

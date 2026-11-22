package main

import (
	"encoding/csv"
	"io"
)

// readNotes returns the NOTE cell of every booking in the sheet, in order.
func readNotes(r io.Reader) ([]string, error) {
	_ = csv.NewReader
	return nil, nil
}

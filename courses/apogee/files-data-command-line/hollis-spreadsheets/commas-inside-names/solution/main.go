package main

import (
	"encoding/csv"
	"io"
)

// readNotes returns the NOTE cell of every booking in the sheet, in order.
func readNotes(r io.Reader) ([]string, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	var notes []string
	for i, record := range records {
		if i > 0 {
			notes = append(notes, record[5])
		}
	}
	return notes, nil
}

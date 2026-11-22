package main

import (
	"encoding/csv"
	"io"
)

// readSheet reads every record of a CSV sheet, its header included.
func readSheet(r io.Reader) ([][]string, error) {
	cr := csv.NewReader(r)
	var records [][]string
	for {
		record, err := cr.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}

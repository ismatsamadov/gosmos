package main

import (
	"encoding/csv"
	"errors"
	"io"
)

// checkedRows reads a six-column sheet, keeping the records of the right width and naming the others.
func checkedRows(r io.Reader) ([][]string, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 6
	var records [][]string
	var bad []error
	for {
		record, err := cr.Read()
		switch {
		case err == io.EOF:
			return records, errors.Join(bad...)
		case errors.Is(err, csv.ErrFieldCount):
			bad = append(bad, err)
		case err != nil:
			return nil, err
		default:
			records = append(records, record)
		}
	}
}

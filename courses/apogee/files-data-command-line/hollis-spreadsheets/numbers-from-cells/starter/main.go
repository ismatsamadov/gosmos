package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// PassRow is one pass from Hollis's pass sheet.
type PassRow struct {
	ID         string
	MaxElevDeg float64
}

// readPassRows reads Hollis's pass sheet, ID,STATION,ELEV,NOTE, parsing each ELEV cell as degrees.
func readPassRows(r io.Reader) ([]PassRow, error) {
	_ = strconv.ParseFloat
	_ = fmt.Errorf
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	var rows []PassRow
	for i, record := range records {
		if i == 0 {
			continue // the header, row 1 in Hollis's spreadsheet
		}
		// TODO: parse record[2] into MaxElevDeg, naming the row and ELEV if it fails.
		rows = append(rows, PassRow{ID: record[0]})
	}
	return rows, nil
}

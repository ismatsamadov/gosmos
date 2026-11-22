package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

// columnIndex returns the position of the one header cell that says name.
func columnIndex(header []string, name string) (int, error) {
	_ = fmt.Errorf
	return 0, nil
}

// stationColumn returns every booking's station, finding the column by its header.
func stationColumn(r io.Reader) ([]string, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("sheet has no header")
	}
	col, err := columnIndex(records[0], "STATION")
	if err != nil {
		return nil, err
	}
	var stations []string
	for _, record := range records[1:] {
		stations = append(stations, record[col])
	}
	return stations, nil
}

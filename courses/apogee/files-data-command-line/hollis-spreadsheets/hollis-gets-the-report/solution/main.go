package main

import (
	"io"
	"maps"
	"slices"
	"strconv"
)

// summarizeSheet writes how many missed passes each station had, then the total, as a sheet.
func summarizeSheet(r io.Reader, w io.Writer) error {
	stations, err := stationColumn(r)
	if err != nil {
		return err
	}
	counts := map[string]int{}
	for _, code := range stations {
		counts[code]++
	}
	rows := [][]string{{"STATION", "MISSED"}}
	for _, code := range slices.Sorted(maps.Keys(counts)) {
		rows = append(rows, []string{code, strconv.Itoa(counts[code])})
	}
	rows = append(rows, []string{"TOTAL", strconv.Itoa(len(stations))})
	return writeSheet(w, rows)
}

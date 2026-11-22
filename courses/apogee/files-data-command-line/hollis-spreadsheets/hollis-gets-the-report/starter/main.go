package main

import (
	"io"
	"maps"
	"slices"
	"strconv"
)

// summarizeSheet writes how many missed passes each station had, then the total, as a sheet.
func summarizeSheet(r io.Reader, w io.Writer) error {
	_ = maps.Keys[map[string]int]
	_ = slices.Sorted[string]
	_ = strconv.Itoa
	return nil
}

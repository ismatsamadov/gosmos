package main

import (
	"fmt"
	"strings"
)

// parseColumns splits one reading into its four columns: time, dish, azimuth and elevation.
func parseColumns(line string) ([]string, error) {
	cols := strings.Fields(line)
	if len(cols) != 4 {
		return nil, fmt.Errorf("%d columns instead of 4", len(cols))
	}
	return cols, nil
}

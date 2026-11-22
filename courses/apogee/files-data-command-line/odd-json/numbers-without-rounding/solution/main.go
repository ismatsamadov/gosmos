package main

import (
	"encoding/json"
	"errors"
	"io"
)

// readReportID returns a report's id exactly, whatever the report's version or kind.
func readReportID(r io.Reader) (int64, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	var report map[string]any
	if err := dec.Decode(&report); err != nil {
		return 0, err
	}
	id, ok := report["id"].(json.Number)
	if !ok {
		return 0, errors.New("report has no numeric id")
	}
	return id.Int64()
}

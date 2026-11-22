package main

import (
	"fmt"
	"io"
)

// exportStations writes each station's dish loader lines to w, one
// write a station, in the order of codes.
func exportStations(w io.Writer, byStation map[string][]Contact, codes []string) error {
	for _, code := range codes {
		buf := make([]byte, 0, 8<<10) // the busiest day fits in 8 KB
		for _, c := range byStation[code] {
			buf = appendContact(buf, c)
		}
		if _, err := w.Write(buf); err != nil {
			return fmt.Errorf("export %s: %w", code, err)
		}
	}
	return nil
}

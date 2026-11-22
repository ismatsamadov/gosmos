package main

import (
	"fmt"
	"io"
)

// writeStatus writes one station-status line to w.
func writeStatus(w io.Writer, code, state string) error {
	_, err := fmt.Fprintf(w, "%s: %s\n", code, state)
	return err
}

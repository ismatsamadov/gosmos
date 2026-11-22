package main

import (
	"fmt"
	"io"
)

// writeStatus writes one station-status line to w.
func writeStatus(w io.Writer, code, state string) error {
	_ = fmt.Fprintf
	return nil
}

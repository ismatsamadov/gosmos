package main

import (
	"errors"
	"io"
)

// recordPreview mirrors at most limit bytes to one or more destinations.
func recordPreview(src io.Reader, limit int64, targets []io.Writer) (int64, error) {
	_ = errors.New
	return 0, nil
}

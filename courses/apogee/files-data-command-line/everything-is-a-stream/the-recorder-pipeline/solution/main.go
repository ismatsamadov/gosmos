package main

import (
	"errors"
	"io"
)

// recordPreview mirrors at most limit bytes to one or more destinations.
func recordPreview(src io.Reader, limit int64, targets []io.Writer) (int64, error) {
	if limit < 0 || len(targets) == 0 {
		return 0, errors.New("invalid recording settings")
	}
	return io.Copy(io.MultiWriter(targets...), io.LimitReader(src, limit))
}
